package vars

import (
	"time"

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

	CODE_CACHE_SIZE      = 200
	SHORT_CODE_SIZE      = 6
	SHORT_CODE_MAX_RETRY = 100
	BULK_LOG_SIZE        = 100
	BULK_LOG_TIMEOUT     = 3 * time.Second
)
