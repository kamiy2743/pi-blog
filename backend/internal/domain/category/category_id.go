package category

import (
	"errors"
	"fmt"
)

type CategoryID string

var errInvalidCategoryID = errors.New("カテゴリIDが不正です")

func ParseCategoryID(s string) (CategoryID, error) {
	if s == "" {
		return "", fmt.Errorf("%w: %q", errInvalidCategoryID, s)
	}
	return CategoryID(s), nil
}
