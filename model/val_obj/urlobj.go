package val_obj

import (
	"database/sql"
	"time"
)

type URLObject struct {
	Id         uint         `json:"id,omitempty"`
	LongURL    string       `json:"long_url,omitempty"`
	ShortCode  string       `json:"short_code,omitempty"`
	CreateTime time.Time    `json:"create_time,omitempty"`
	ExpireTime sql.NullTime `json:"expire_time,omitempty"`
	Enabled    bool         `json:"enabled,omitempty"`
}
