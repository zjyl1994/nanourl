package render_obj

type URLObject struct {
	Id         uint   `json:"id"`
	LongURL    string `json:"long_url"`
	ShortCode  string `json:"short_code"`
	CreateTime int64  `json:"create_time"`
	ExpireTime int64  `json:"expire_time"`
	Enabled    bool   `json:"enabled"`
	ClickCount int    `json:"click_count"`
	HrefLink   string `json:"href_link"`
}
