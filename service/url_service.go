package service

import (
	"errors"

	"github.com/zjyl1994/nanourl/model"
	"github.com/zjyl1994/nanourl/util"
	"github.com/zjyl1994/nanourl/vars"
	"gorm.io/gorm"
)

type URLService struct{}

func (URLService) New(target, code string) (uint, error) {
	var m model.URLObject
	m.URL = target
	m.Code = code
	err := vars.DB.Create(&m).Error
	if err != nil {
		if util.IsSqliteDuplicateError(err) {
			return 0, ErrDuplicate
		}
		return 0, err
	}
	return m.ID, nil
}

func (URLService) Get(id uint) (*model.URLObject, error) {
	var m model.URLObject
	err := vars.DB.First(&m, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &m, nil
}

func (URLService) List(page, pageSize int) ([]model.URLObject, int64, error) {
	var totalCount int64
	var datas []model.URLObject

	err := vars.DB.Model(&model.URLObject{}).Count(&totalCount).Error
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
