package formatter

import "time"

func FormatTimeISO8601(value *time.Time) string {
	if value == nil {
		return ""
	}
	return value.Format(time.RFC3339)
}
