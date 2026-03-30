package category

import "context"

type CategoryRepository interface {
	Create(ctx context.Context, input CreateCategoryInput) (Category, error)
	Update(ctx context.Context, category Category) error
	All(ctx context.Context) ([]Category, error)
	Search(ctx context.Context, criteria SearchCategoryCriteria) ([]Category, error)
}
