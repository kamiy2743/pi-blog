package domain

type OrderDirection string

const (
	OrderDirectionAsc  OrderDirection = "asc"
	OrderDirectionDesc OrderDirection = "desc"
)

type OrderBy struct {
	Column    string
	Direction OrderDirection
}
