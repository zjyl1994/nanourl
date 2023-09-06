package vars

import (
	"github.com/speps/go-hashids/v2"
	"gorm.io/gorm"
)

var (
	Listen  string
	DataDir string
	BaseUrl string

	DB     *gorm.DB
	HashId *hashids.HashID
)

const (
	DEFAULT_LISTEN  = ":9900"
	DEFAULT_BASEURL = "http://localhost:9900/"
)
