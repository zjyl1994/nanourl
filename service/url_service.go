package service

import (
	"errors"

	"github.com/zjyl1994/nanourl/model/db_model"
	"github.com/zjyl1994/nanourl/model/val_obj"
	"github.com/zjyl1994/nanourl/util"
	"github.com/zjyl1994/nanourl/vars"
	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"
)

var getUrlObjectFlight singleflight.Group

type URLService struct{}

func (URLService) New(val val_obj.URLObject) (uint, error) {
	var m db_model.URLObject
	m.URL = val.LongURL
	m.Code = val.ShortCode
	err := vars.DB.Create(&m).Error
	if err != nil {
		if util.IsSqliteDuplicateError(err) {
			return 0, ErrDuplicate
		}
		return 0, err
	}
	return m.ID, nil
}

func (URLService) GetById(id uint) (*db_model.URLObject, error) {
	var m db_model.URLObject
	err := vars.DB.First(&m, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &m, nil
}

func (URLService) List(page, pageSize int) ([]db_model.URLObject, int64, error) {
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
	return datas, totalCount, nil
}

func (URLService) GetByCode(code string) (*db_model.URLObject, error) {
	var m db_model.URLObject
	err := vars.DB.First(&m, "code = ?", code).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &m, nil
}

func (URLService) GetShortCode(id uint) (string, error) {
	return vars.HashId.Encode([]int{int(id)})
}

func (s URLService) SearchCode(code string) (*db_model.URLObject, error) {
	if obj, err := s.GetByCode(code); err == nil {
		return obj, nil
	}
	ids, err := vars.HashId.DecodeWithError(code)
	if err != nil {
		return nil, ErrNotFound
	}
	return s.GetById(uint(ids[0]))
}

func (s URLService) SearchCodeWithCache(code string) (val_obj.URLObject, error) {
	if obj, ok := vars.CodeCache.Get(code); ok {
		return obj, nil
	} else {
		result, err, _ := getUrlObjectFlight.Do(code, func() (interface{}, error) {
			return s.SearchCode(code)
		})
		if err != nil {
			return val_obj.URLObject{}, nil
		}
		urlObj := result.(*db_model.URLObject)
		val := val_obj.URLObject{
			Id:        urlObj.ID,
			LongURL:   urlObj.URL,
			ShortCode: urlObj.Code,
		}
		vars.CodeCache.Add(code, val)
		return val, nil
	}
}
