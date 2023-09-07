package vars

import (
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/zjyl1994/nanourl/model/val_obj"
	"gorm.io/gorm"
)

var (
	Listen       string
	DataDir      string
	BaseUrl      string
	RealIpHeader string

	DB        *gorm.DB
	CodeCache *lru.TwoQueueCache[string, val_obj.URLObject]
)

const (
	DEFAULT_LISTEN  = ":9900"
	DEFAULT_BASEURL = "http://localhost:9900/"

	CODE_CACHE_SIZE = 200
	SHORT_URL_SIZE  = 6
)
