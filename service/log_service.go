package service

import (
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
	geodb, err := geoip2.Open(filepath.Join(vars.DataDir, vars.GEOIP_DATABASE_FILENAME))
	if err == nil {
		defer geodb.Close()
	} else {
		log.Errorln(err.Error())
	}
	for {
		items, length, _, ok := lo.BufferWithTimeout(backgroundLogChan, vars.BULK_LOG_SIZE, vars.BULK_LOG_TIMEOUT)
		if length == 0 {
			continue
		}
		models := lo.Map(items, func(v val_obj.AccessLog, _ int) db_model.AccessLog {
			return db_model.AccessLog{
				UrlId:       v.UrlId,
				Referrer:    v.Referrer,
				UserIp:      v.UserIp,
				UserAgent:   v.UserAgent,
				UserCountry: getIPCountry(geodb, v.UserIp),
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

func getIPCountry(geodb *geoip2.Reader, ip string) string {
	if geodb != nil {
		ipAddr := net.ParseIP(ip)
		if ipAddr != nil {
			r, err := geodb.Country(ipAddr)
			if err == nil {
				return r.Country.IsoCode
			}
		}
	}
	return ""
}

func (LogService) List(urlId, page, pageSize int) ([]val_obj.AccessLog, int64, error) {
	var totalCount int64
	var datas []db_model.AccessLog

	countQuery := vars.DB.Model(&db_model.AccessLog{})
	if urlId > 0 {
		countQuery = countQuery.Where("url_id = ?", urlId)
	}
	err := countQuery.Count(&totalCount).Error
	if err != nil {
		return nil, 0, err
	}

	limit, offset := util.PageHelper(page, pageSize)
	dataQuery := vars.DB.Limit(limit).Offset(offset)
	if urlId > 0 {
		dataQuery = dataQuery.Where("url_id = ?", urlId)
	}
	err = dataQuery.Find(&datas).Error
	if err != nil {
		return nil, 0, err
	}

	var results []val_obj.AccessLog
	for _, v := range datas {
		var country string
		if v.UserCountry != "" {
			if c, ok := vars.GeoCountry[v.UserCountry]; ok {
				country = c.Emoji + " " + c.Name
			} else {
				country = v.UserCountry
			}
		}
		results = append(results, val_obj.AccessLog{
			UrlId:       v.UrlId,
			Referrer:    v.Referrer,
			UserIp:      v.UserIp,
			UserCountry: country,
			UserAgent:   v.UserAgent,
			AccessTime:  v.CreatedAt,
		})
	}

	return results, totalCount, nil
}

func (LogService) CountLog(urlid []int) (map[int]int, error) {
	type resultContainer struct {
		UrlId int `gorm:"column:url_id"`
		Total int `gorm:"total"`
	}
	var result []resultContainer
	err := vars.DB.Model(&db_model.AccessLog{}).Select("url_id,count(*) as total").Where("url_id IN (?)", urlid).Group("url_id").Find(&result).Error
	if err != nil {
		return nil, err
	}
	return lo.SliceToMap(result, func(x resultContainer) (int, int) { return x.UrlId, x.Total }), nil
}
