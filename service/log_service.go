package service

import (
	"github.com/zjyl1994/nanourl/model/db_model"
	"github.com/zjyl1994/nanourl/model/val_obj"
	"github.com/zjyl1994/nanourl/vars"
)

type LogService struct{}

func (LogService) AddLog(val val_obj.AccessLog) error {
	var m db_model.AccessLog
	m.UrlId = val.UrlId
	m.Referrer = val.Referrer
	m.UserAgent = val.UserAgent
	m.UserIp = val.UserIp

	return vars.DB.Create(&m).Error
}
