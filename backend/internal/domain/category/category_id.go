package category

import (
	"fmt"

	"blog/internal/domain"
)

type CategoryID uint32

func ParseCategoryID(s string) (CategoryID, error) {
	id, err := domain.ParseUint32(s)
	if err != nil {
		return 0, fmt.Errorf("カテゴリIDが不正です: %q", s)
	}
	return CategoryID(id), nil
}
