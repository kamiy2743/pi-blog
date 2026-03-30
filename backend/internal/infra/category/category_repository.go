package category

import (
	"context"

	"blog/internal/domain/category"
)

type CategoryRepository struct{}

func NewCategoryRepository() *CategoryRepository {
	return &CategoryRepository{}
}

func (r *CategoryRepository) All(_ context.Context) ([]category.Category, error) {
	raspiCategoryID, _ := category.ParseCategoryID("1")
	infraCategoryID, _ := category.ParseCategoryID("2")
	backendCategoryID, _ := category.ParseCategoryID("3")

	return []category.Category{
		{ID: raspiCategoryID, Name: "Raspberry Pi"},
		{ID: infraCategoryID, Name: "Infrastructure"},
		{ID: backendCategoryID, Name: "Backend"},
	}, nil
}

func (r *CategoryRepository) Create(_ context.Context, input category.CreateCategoryInput) (category.Category, error) {
	if err := input.Validate(); err != nil {
		return category.Category{}, err
	}

	categoryID, _ := category.ParseCategoryID("1")
	return category.Category{
		ID:   categoryID,
		Name: input.Name,
	}, nil
}

func (r *CategoryRepository) Search(_ context.Context, criteria category.SearchCategoryCriteria) ([]category.Category, error) {
	categories, err := r.All(context.Background())
	if err != nil {
		return nil, err
	}

	filtered := make([]category.Category, 0, len(categories))
	for _, category := range categories {
		if !matchesCategoryIDs(criteria.IDs, category.ID) {
			continue
		}
		filtered = append(filtered, category)
	}

	return filtered, nil
}

func (r *CategoryRepository) Update(_ context.Context, category category.Category) error {
	return category.Validate()
}

func matchesCategoryIDs(ids []category.CategoryID, target category.CategoryID) bool {
	if len(ids) == 0 {
		return true
	}

	for _, id := range ids {
		if id == target {
			return true
		}
	}
	return false
}
