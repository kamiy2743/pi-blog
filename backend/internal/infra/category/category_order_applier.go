package category

import (
	"fmt"

	"blog/internal/db/ent"
	entCategory "blog/internal/db/ent/category"
	domainCategory "blog/internal/domain/category"
)

func ApplyOrder(query *ent.CategoryQuery, orderBy domainCategory.OrderBy) error {
	if orderBy == "" {
		orderBy = domainCategory.OrderByDefault
	}
	switch orderBy {
	case domainCategory.OrderByNameAsc:
		query.Order(entCategory.ByName())
		return nil
	default:
		return fmt.Errorf("未対応のカテゴリの並び順です: %s", orderBy)
	}
}
