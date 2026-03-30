package article

import "blog/internal/domain"

type SearchArticleCriteria struct {
	Title   string
	Limit   int
	OrderBy domain.OrderBy
}
