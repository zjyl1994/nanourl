package render_obj

type AccessLog struct {
	UrlId        uint   `json:"url_id"`
	Referrer     string `json:"referrer"`
	UserIp       string `json:"user_ip"`
	CountryCode  string `json:"country_code"`
	CountryName  string `json:"country_name"`
	CountryEmoji string `json:"country_emoji"`
	UserAgent    string `json:"user_agent"`
	AccessTime   int64  `json:"access_time"`
}
