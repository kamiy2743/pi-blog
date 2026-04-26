package domain

import (
	"errors"
	"strconv"
)

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
