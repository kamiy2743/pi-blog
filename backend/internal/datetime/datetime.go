package datetime

import (
	"errors"
	"time"
)

func Now() time.Time {
	return time.Now().UTC()
}

func Parse(value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, errors.New("文字列が空です")
	}
	parsed, err := time.Parse(time.RFC3339Nano, value)
	if err != nil {
		return time.Time{}, err
	}
	return parsed.UTC(), nil
}

func ParseOptional(value string) (*time.Time, error) {
	if value == "" {
		return nil, nil
	}

	parsed, err := Parse(value)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

func FormatISO8601(value *time.Time) string {
	if value == nil {
		return ""
	}
	return value.UTC().Format(time.RFC3339)
}
