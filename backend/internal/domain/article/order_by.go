package article

type OrderBy string

const (
	OrderByDefault OrderBy = OrderByLatest
	OrderByLatest  OrderBy = "latest"
)
