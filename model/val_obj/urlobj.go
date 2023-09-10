package val_obj

import "time"

type URLObject struct {
	Id         uint      `json:"id,omitempty"`
	LongURL    string    `json:"long_url,omitempty"`
	ShortCode  string    `json:"short_code,omitempty"`
	CreateTime time.Time `json:"create_time,omitempty"`
}
