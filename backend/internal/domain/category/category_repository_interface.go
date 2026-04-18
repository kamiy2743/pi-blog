package category

import (
	"context"
)

type CategoryRepository interface {
	Create(ctx context.Context, input CreateCategoryInput) error
	Update(ctx context.Context, category Category) error
	Delete(ctx context.Context, category Category) error
	All(ctx context.Context, orderBy OrderBy) ([]Category, error)
	Search(ctx context.Context, criteria SearchCategoryCriteria) ([]Category, error)
}
