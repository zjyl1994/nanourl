package service

import (
	"errors"

	"github.com/samber/lo"
	"github.com/zjyl1994/nanourl/model/db_model"
	"github.com/zjyl1994/nanourl/model/val_obj"
	"github.com/zjyl1994/nanourl/util"
	"github.com/zjyl1994/nanourl/vars"
	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"
)

var getUrlObjectFlight singleflight.Group

type URLService struct{}

func (URLService) New(val val_obj.URLObject) (uint, string, error) {
	var retry int
	for retry < vars.SHORT_CODE_MAX_RETRY {
		var m db_model.URLObject
		m.URL = val.LongURL

		if val.ShortCode != "" {
			m.Code = val.ShortCode
		} else {
			m.Code = util.RandString(vars.SHORT_CODE_SIZE)
		}

		err := vars.DB.Create(&m).Error
		if err == nil {
			return m.ID, m.Code, nil
		}
		if !util.IsSqliteDuplicateError(err) {
			return 0, "", err

		}
		if val.ShortCode != "" {
			return 0, "", ErrCodeDuplicate
		}
		retry++
	}
	return 0, "", ErrCodeExhausted
}

func (URLService) GetById(id uint) (*val_obj.URLObject, error) {
	var m db_model.URLObject
	err := vars.DB.First(&m, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCodeNotFound
		}
		return nil, err
	}
	return &val_obj.URLObject{
		Id:         m.ID,
		LongURL:    m.URL,
		ShortCode:  m.Code,
		CreateTime: m.CreatedAt,
	}, nil
}

func (URLService) List(page, pageSize int) ([]val_obj.URLObject, int64, error) {
	var totalCount int64
	var datas []db_model.URLObject

	err := vars.DB.Model(&db_model.URLObject{}).Count(&totalCount).Error
	if err != nil {
		return nil, 0, err
	}

	limit, offset := util.PageHelper(page, pageSize)
	err = vars.DB.Limit(limit).Offset(offset).Find(&datas).Error
	if err != nil {
		return nil, 0, err
	}
	return lo.Map(datas, func(x db_model.URLObject, _ int) val_obj.URLObject {
		return val_obj.URLObject{
			Id:         x.ID,
			LongURL:    x.URL,
			ShortCode:  x.Code,
			CreateTime: x.CreatedAt,
		}
	}), totalCount, nil
}

func (URLService) GetByCode(code string) (*val_obj.URLObject, error) {
	var m db_model.URLObject
	err := vars.DB.First(&m, "code = ?", code).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCodeNotFound
		}
		return nil, err
	}
	return &val_obj.URLObject{
		Id:         m.ID,
		LongURL:    m.URL,
		ShortCode:  m.Code,
		CreateTime: m.CreatedAt,
	}, nil
}

func (s URLService) SearchCode(code string) (val_obj.URLObject, error) {
	if obj, ok := vars.CodeCache.Get(code); ok {
		return obj, nil
	} else {
		result, err, _ := getUrlObjectFlight.Do(code, func() (interface{}, error) {
			return s.GetByCode(code)
		})
		if err != nil {
			return val_obj.URLObject{}, err
		}
		urlObj := result.(*db_model.URLObject)
		val := val_obj.URLObject{
			Id:         urlObj.ID,
			LongURL:    urlObj.URL,
			ShortCode:  urlObj.Code,
			CreateTime: urlObj.CreatedAt,
		}
		vars.CodeCache.Add(code, val)
		return val, nil
	}
}
