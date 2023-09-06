package model

import "gorm.io/gorm"

type URLObject struct {
	gorm.Model
	URL  string
	Code string `gorm:"uniqueIndex"`
}
