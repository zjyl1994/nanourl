package db_model

import "gorm.io/gorm"

type AccessLog struct {
	gorm.Model
	UrlId       uint
	Referrer    string
	UserIp      string
	UserAgent   string
	UserCountry string
}
