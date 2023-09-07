package db_model

import "gorm.io/gorm"

type URLObject struct {
	gorm.Model
	URL  string
	Code string `gorm:"uniqueIndex"`
}
