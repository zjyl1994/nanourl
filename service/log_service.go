package service

import (
	"sync"

	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"github.com/zjyl1994/nanourl/model/db_model"
	"github.com/zjyl1994/nanourl/model/val_obj"
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
