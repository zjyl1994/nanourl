package val_obj

import "time"

type AccessLog struct {
	UrlId       uint      `json:"url_id,omitempty"`
	Referrer    string    `json:"referrer,omitempty"`
	UserIp      string    `json:"user_ip,omitempty"`
	UserCountry string    `json:"user_country,omitempty"`
	UserAgent   string    `json:"user_agent,omitempty"`
	AccessTime  time.Time `json:"access_time,omitempty"`

	OS         string `json:"os,omitempty"`
	Browser    string `json:"browser,omitempty"`
	Device     string `json:"device,omitempty"`
	DeviceType string `json:"device_type,omitempty"`
}
