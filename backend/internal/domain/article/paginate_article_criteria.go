package article

type PaginateArticleCriteria struct {
	SearchCriteria SearchArticleCriteria
	Page           int
	PerPage        int
}
