package db_model

import (
	"database/sql"

	"gorm.io/gorm"
)

type URLObject struct {
	gorm.Model
	URL      string
	Code     string `gorm:"uniqueIndex"`
	ExpireAt sql.NullTime
	Enabled  bool
}
