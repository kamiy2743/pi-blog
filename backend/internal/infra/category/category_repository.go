package category

import (
	"context"

	"blog/internal/db/ent"
	entCategory "blog/internal/db/ent/category"
	domainCategory "blog/internal/domain/category"
)

type CategoryRepository struct {
	client *ent.Client
}

func NewCategoryRepository(client *ent.Client) *CategoryRepository {
	return &CategoryRepository{client: client}
}

func (r *CategoryRepository) Create(ctx context.Context, input domainCategory.CreateCategoryInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	if _, err := r.client.Category.Create().
		SetName(input.Name).
		Save(ctx); err != nil {
		return err
	}

	return nil
}

func (r *CategoryRepository) Update(ctx context.Context, entity domainCategory.Category) error {
	if err := entity.Validate(); err != nil {
		return err
	}

	return r.client.Category.UpdateOneID(uint32(entity.ID)).
		SetName(entity.Name).
		Exec(ctx)
}

func (r *CategoryRepository) Delete(ctx context.Context, entity domainCategory.Category) error {
	return r.client.Category.DeleteOneID(uint32(entity.ID)).Exec(ctx)
}

func (r *CategoryRepository) All(ctx context.Context, orderBy domainCategory.OrderBy) ([]domainCategory.Category, error) {
	query := r.client.Category.Query()

	if err := ApplyOrder(query, orderBy); err != nil {
		return nil, err
	}

	models, err := query.All(ctx)
	if err != nil {
		return nil, err
	}
	return HydrateCategories(models), nil
}

func (r *CategoryRepository) Search(ctx context.Context, criteria domainCategory.SearchCategoryCriteria) ([]domainCategory.Category, error) {
	query := r.client.Category.Query()

	if len(criteria.IDs) > 0 {
		categoryIDs := make([]uint32, 0, len(criteria.IDs))
		for _, id := range criteria.IDs {
			categoryIDs = append(categoryIDs, uint32(id))
		}
		query.Where(entCategory.IDIn(categoryIDs...))
	}

	if criteria.Name != "" {
		query.Where(entCategory.Name(criteria.Name))
	}

	if err := ApplyOrder(query, criteria.OrderBy); err != nil {
		return nil, err
	}

	models, err := query.All(ctx)
	if err != nil {
		return nil, err
	}
	return HydrateCategories(models), nil
}
