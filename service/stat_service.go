package service

import (
	"fmt"

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
		Joins("LEFT JOIN url_objects ON access_logs.url_id = url_objects.id").
		Where(fmt.Sprintf("access_logs.created_at > DATE('now', '-%d days')", vars.STATISTICS_DAY)).
		Group("url_objects.code").
		Order("click DESC").Limit(vars.STATISTICS_LIMIT).Find(&dbm).Error
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
	err := vars.DB.Table("access_logs").Select("user_country as dim,count(*) as click").
		Where(fmt.Sprintf("created_at > DATE('now', '-%d days')", vars.STATISTICS_DAY)).
		Group("dim").Order("click DESC").Limit(vars.STATISTICS_LIMIT).Find(&dbm).Error
	if err != nil {
		return nil, err
	}
	return lo.Map(dbm, func(x db_model.StatItem, _ int) val_obj.StatItem {
		country := util.CountryCode2EmojiAndName(x.Dim)
		return val_obj.StatItem{Dim: country, Click: x.Click}
	}), nil
}

func (StatService) DayClick() ([]val_obj.StatItem, error) {
	var dbm []db_model.StatItem
	err := vars.DB.Table("access_logs").Select("DATE(created_at) as dim,count(*) as click").
		Where(fmt.Sprintf("created_at > DATE('now', '-%d days')", vars.STATISTICS_DAY)).
		Group("dim").Order("dim DESC").Limit(vars.STATISTICS_LIMIT).Find(&dbm).Error
	if err != nil {
		return nil, err
	}
	return lo.Map(dbm, func(x db_model.StatItem, _ int) val_obj.StatItem {
		return val_obj.StatItem{Dim: x.Dim, Click: x.Click}
	}), nil
}

func (s StatService) TopOS() ([]val_obj.StatItem, error) {
	return s.fieldCounter("os")
}

func (s StatService) TopBrowser() ([]val_obj.StatItem, error) {
	return s.fieldCounter("browser")
}

func (s StatService) TopDevice() ([]val_obj.StatItem, error) {
	return s.fieldCounter("device")
}

func (s StatService) fieldCounter(fieldName string) ([]val_obj.StatItem, error) {
	var dbm []db_model.StatItem
	err := vars.DB.Table("access_logs").Select(fieldName + " as dim,count(*) as click").
		Where(fmt.Sprintf("created_at > DATE('now', '-%d days')", vars.STATISTICS_DAY)).
		Group("dim").Order("dim DESC").Limit(vars.STATISTICS_LIMIT).Find(&dbm).Error
	if err != nil {
		return nil, err
	}
	return lo.Map(dbm, func(x db_model.StatItem, _ int) val_obj.StatItem {
		var dim string
		if x.Dim == "" {
			dim = "Unknown"
		} else {
			dim = x.Dim
		}
		return val_obj.StatItem{Dim: dim, Click: x.Click}
	}), nil
}
