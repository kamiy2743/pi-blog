package category

import (
	"context"
	"fmt"

	domainCategory "blog/internal/domain/category"
	"blog/internal/ent"
	entCategory "blog/internal/ent/category"
)

type CategoryRepository struct {
	client *ent.Client
}

func NewCategoryRepository(client *ent.Client) *CategoryRepository {
	return &CategoryRepository{client: client}
}

func (r *CategoryRepository) Create(ctx context.Context, input domainCategory.CreateCategoryInput) (domainCategory.Category, error) {
	return domainCategory.Category{}, nil
}

func (r *CategoryRepository) Update(ctx context.Context, input domainCategory.Category) error {
	return nil
}

func (r *CategoryRepository) All(ctx context.Context, orderBy domainCategory.OrderBy) ([]domainCategory.Category, error) {
	query := r.client.Category.Query()

	switch orderBy {
	case domainCategory.OrderByNameAsc:
		query.Order(entCategory.ByName())
	default:
		return nil, fmt.Errorf("未対応のカテゴリの並び順です: %s", orderBy)
	}

	models, err := query.All(ctx)
	if err != nil {
		return nil, err
	}
	return HydrateCategories(models), nil
}

func (r *CategoryRepository) Search(ctx context.Context, criteria domainCategory.SearchCategoryCriteria) ([]domainCategory.Category, error) {
	return nil, nil
}
