package category

import (
	"errors"
	"fmt"

	"blog/internal/domain"
)

type CategoryID uint32

var errInvalidCategoryID = errors.New("カテゴリIDが不正です")

func ParseCategoryID(s string) (CategoryID, error) {
	id, err := domain.ParseUint32(s)
	if err != nil {
		return 0, fmt.Errorf("%w: %q", errInvalidCategoryID, s)
	}
	return CategoryID(id), nil
}
