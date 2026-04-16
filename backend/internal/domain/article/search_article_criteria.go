package article

import "blog/internal/domain/category"

type SearchArticleCriteria struct {
	Title              string
	CategoryIDs        []category.CategoryID
	IncludeUnpublished bool
	Limit              *int
	OrderBy            OrderBy
}
