package val_obj

type URLObject struct {
	Id        uint   `json:"id,omitempty"`
	LongURL   string `json:"long_url,omitempty"`
	ShortCode string `json:"short_code,omitempty"`
}
