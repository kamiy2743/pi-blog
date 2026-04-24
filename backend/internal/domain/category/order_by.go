package category

type OrderBy string

const (
	OrderByDefault OrderBy = OrderByNameAsc
	OrderByNameAsc OrderBy = "name_asc"
)
