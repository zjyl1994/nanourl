package val_obj

import "time"

type AccessLog struct {
	UrlId        uint
	Referrer     string
	UserIp       string
	UserLocation string
	UserAgent    string
	AccessTime   time.Time
}
