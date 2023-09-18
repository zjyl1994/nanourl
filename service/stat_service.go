package service

import (
	"github.com/samber/lo"
	"github.com/zjyl1994/nanourl/model/db_model"
	"github.com/zjyl1994/nanourl/model/val_obj"
	"github.com/zjyl1994/nanourl/util"
	"github.com/zjyl1994/nanourl/vars"
)

type StatService struct{}

func (StatService) TopURL() ([]val_obj.StatItem, error) {
	var dbm []db_model.StatItem
	err := vars.DB.Table("access_logs").Select("url_objects.code as dim,count(*) as click").
		Joins("LEFT JOIN url_objects ON access_logs.url_id = url_objects.id").Group("url_objects.code").
		Order("click DESC").Limit(10).Find(&dbm).Error
	if err != nil {
		return nil, err
	}
	return lo.Map(dbm, func(x db_model.StatItem, _ int) val_obj.StatItem {
		url := vars.BaseUrl + x.Dim
		return val_obj.StatItem{Dim: url, Click: x.Click}
	}), nil
}

func (StatService) TopCountry() ([]val_obj.StatItem, error) {
	var dbm []db_model.StatItem
	err := vars.DB.Table("access_logs").Select("access_logs.user_country as dim,count(*) as click").
		Group("access_logs.user_country").
		Order("click DESC").Limit(10).Find(&dbm).Error
	if err != nil {
		return nil, err
	}
	return lo.Map(dbm, func(x db_model.StatItem, _ int) val_obj.StatItem {
		country := util.CountryCode2EmojiAndName(x.Dim)
		return val_obj.StatItem{Dim: country, Click: x.Click}
	}), nil
}
