package category

import (
	"context"

	domainCategory "blog/internal/domain/category"
)

type CategoryRepositoryStub struct {
	CreateFunc func(ctx context.Context, input domainCategory.CreateCategoryInput) error
	UpdateFunc func(ctx context.Context, category domainCategory.Category) error
	DeleteFunc func(ctx context.Context, category domainCategory.Category) error
	AllFunc    func(ctx context.Context, orderBy domainCategory.OrderBy) ([]domainCategory.Category, error)
	SearchFunc func(ctx context.Context, criteria domainCategory.SearchCategoryCriteria) ([]domainCategory.Category, error)
}

func (s CategoryRepositoryStub) Create(ctx context.Context, input domainCategory.CreateCategoryInput) error {
	if s.CreateFunc == nil {
		return nil
	}
	return s.CreateFunc(ctx, input)
}

func (s CategoryRepositoryStub) Update(ctx context.Context, category domainCategory.Category) error {
	if s.UpdateFunc == nil {
		return nil
	}
	return s.UpdateFunc(ctx, category)
}

func (s CategoryRepositoryStub) Delete(ctx context.Context, category domainCategory.Category) error {
	if s.DeleteFunc == nil {
		return nil
	}
	return s.DeleteFunc(ctx, category)
}

func (s CategoryRepositoryStub) All(ctx context.Context, orderBy domainCategory.OrderBy) ([]domainCategory.Category, error) {
	if s.AllFunc == nil {
		return nil, nil
	}
	return s.AllFunc(ctx, orderBy)
}

func (s CategoryRepositoryStub) Search(ctx context.Context, criteria domainCategory.SearchCategoryCriteria) ([]domainCategory.Category, error) {
	if s.SearchFunc == nil {
		return nil, nil
	}
	return s.SearchFunc(ctx, criteria)
}
