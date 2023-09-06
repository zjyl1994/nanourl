package vars

import "gorm.io/gorm"

var (
	Listen  string
	DataDir string

	DB *gorm.DB
)
