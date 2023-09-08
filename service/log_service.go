package service

import (
	"fmt"
	"net"
	"path/filepath"
	"sync"

	"github.com/oschwald/geoip2-golang"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"github.com/zjyl1994/nanourl/model/db_model"
	"github.com/zjyl1994/nanourl/model/val_obj"
	"github.com/zjyl1994/nanourl/util"
	"github.com/zjyl1994/nanourl/vars"
)

var backgroundLogOnce sync.Once
var backgroundLogChan = make(chan val_obj.AccessLog, vars.BULK_LOG_SIZE*2)

type LogService struct{}

func (s LogService) AddLog(val val_obj.AccessLog) {
	backgroundLogOnce.Do(func() {
		go s.backgroundLogWorker()
	})

	backgroundLogChan <- val
}

func (LogService) backgroundLogWorker() {
	for {
		items, length, _, ok := lo.BufferWithTimeout(backgroundLogChan, vars.BULK_LOG_SIZE, vars.BULK_LOG_TIMEOUT)
		if length == 0 {
			continue
		}
		models := lo.Map(items, func(v val_obj.AccessLog, _ int) db_model.AccessLog {
			return db_model.AccessLog{
				UrlId:     v.UrlId,
				Referrer:  v.Referrer,
				UserIp:    v.UserIp,
				UserAgent: v.UserAgent,
			}
		})
		err := vars.DB.CreateInBatches(models, vars.BULK_LOG_SIZE).Error
		if err != nil {
			log.Errorln("bulk write access_log error", err.Error())
		}
		if !ok {
			break
		}
	}
}

func (LogService) List(page, pageSize int) ([]val_obj.AccessLog, int64, error) {
	var totalCount int64
	var datas []db_model.AccessLog

	err := vars.DB.Model(&db_model.AccessLog{}).Count(&totalCount).Error
	if err != nil {
		return nil, 0, err
	}

	limit, offset := util.PageHelper(page, pageSize)
	err = vars.DB.Limit(limit).Offset(offset).Find(&datas).Error
	if err != nil {
		return nil, 0, err
	}

	results, err := dbAccessLog2valobj(datas)
	if err != nil {
		return nil, 0, err
	}

	return results, totalCount, nil
}

func dbAccessLog2valobj(inputs []db_model.AccessLog) ([]val_obj.AccessLog, error) {
	geodb, err := geoip2.Open(filepath.Join(vars.DataDir, vars.GEOIP_DATABASE_FILENAME))
	if err != nil {
		return nil, err
	}
	defer geodb.Close()
	result := make([]val_obj.AccessLog, 0)
	for _, v := range inputs {
		var userLocation string
		if v.UserIp != "" {
			if ip := net.ParseIP(v.UserIp); ip != nil {
				if r, err := geodb.Country(ip); err == nil {
					code := r.Country.IsoCode
					if code != "" {
						name := r.Country.Names["en"]
						emoji := string('🇦'+rune(code[0])-'A') + string('🇦'+rune(code[1])-'A')
						userLocation = fmt.Sprintf("%s %s (%s)", emoji, name, code)
					}
				}
			}
		}
		result = append(result, val_obj.AccessLog{
			UrlId:        v.UrlId,
			Referrer:     v.Referrer,
			UserIp:       v.UserIp,
			UserLocation: userLocation,
			UserAgent:    v.UserAgent,
			AccessTime:   v.CreatedAt,
		})
	}
	return result, nil
}
