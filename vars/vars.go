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
	HomepageUrl  string

	DB         *gorm.DB
	CodeCache  *lru.TwoQueueCache[string, val_obj.URLObject]
	GeoCountry map[string]GeoCountryItem
)

const (
	DEFAULT_LISTEN  = ":9900"
	DEFAULT_BASEURL = "http://localhost:9900/"

	CODE_CACHE_SIZE          = 200
	SHORT_CODE_SIZE          = 6
	SHORT_CODE_MAX_RETRY     = 100
	BULK_LOG_SIZE            = 100
	BULK_LOG_TIMEOUT         = 3 * time.Second
	DEFAULT_DOWNLOAD_TIMEOUT = 10 * time.Second

	GEOIP_DOWNLOAD_URL      = "https://raw.githubusercontent.com/P3TERX/GeoLite.mmdb/download/GeoLite2-Country.mmdb"
	GEOIP_DATABASE_FILENAME = "GeoLite2-Country.mmdb"
	GEOIP_EMOJI_URL         = "https://cdn.jsdelivr.net/npm/country-flag-emoji-json@2.0.0/dist/index.json"
	GEOIP_EMOJI_FILENAME    = "flag-emoji.json"
)

type GeoCountryItem struct {
	Name    string `json:"name"`
	Code    string `json:"code"`
	Emoji   string `json:"emoji"`
	Unicode string `json:"unicode"`
	Image   string `json:"image"`
}
