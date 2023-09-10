package util

import (
	"database/sql"
	"time"
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
