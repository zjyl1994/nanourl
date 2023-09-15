package util

import (
	"database/sql"
	"time"

	"github.com/mileusna/useragent"
)

func FormatTime(t time.Time) string {
	return t.Format(time.DateTime)
}

func FormatNullableTime(t sql.NullTime) string {
	if t.Valid {
		return t.Time.Format(time.DateTime)
	}
	return ""
}

type UserAgentData struct {
	OS         string
	Browser    string
	Device     string
	DeviceType string
}

func ParseUserAgent(s string) UserAgentData {
	var result UserAgentData
	ua := useragent.Parse(s)

	result.OS = ua.OS + " " + ua.OSVersion
	result.Browser = ua.Name + " " + ua.Version
	result.Device = ua.Device
	switch {
	case ua.Desktop:
		result.DeviceType = "Desktop"
	case ua.Tablet:
		result.DeviceType = "Tablet"
	case ua.Mobile:
		result.DeviceType = "Mobile"
	case ua.Bot:
		result.DeviceType = "Bot"
	}
	return result
}
