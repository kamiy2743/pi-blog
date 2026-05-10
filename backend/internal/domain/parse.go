package domain

import (
	"errors"
	"strconv"
	"time"
)

const DatetimeLayout = "2006-01-02T15:04"

func ParseBool(s string) (bool, error) {
	return strconv.ParseBool(s)
}

func ParseDatetime(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, errors.New("文字列が空です")
	}

	parsed, err := time.Parse(DatetimeLayout, s)
	if err != nil {
		return time.Time{}, err
	}

	return parsed, nil
}

func ParseOptionalDatetime(s string) (*time.Time, error) {
	if s == "" {
		return nil, nil
	}

	parsed, err := ParseDatetime(s)
	if err != nil {
		return nil, err
	}

	return &parsed, nil
}

func ParseUint32(s string) (uint32, error) {
	if s == "" {
		return 0, errors.New("文字列が空です")
	}

	id, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, err
	}

	return uint32(id), nil
}
