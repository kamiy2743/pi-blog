package article

import "blog/internal/domain/category"

type SearchArticleCriteria struct {
	IDs                []ArticleID
	Title              string
	CategoryIDs        []category.CategoryID
	IncludeUnpublished bool
	Limit              *int
	OrderBy            OrderBy
}
