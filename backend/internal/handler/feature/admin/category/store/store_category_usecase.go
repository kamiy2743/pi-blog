package store

import (
	"context"

	"blog/internal/domain/category"
	"blog/internal/handler/handlererror"
)

type Usecase struct {
	categoryRepository category.CategoryRepository
}

func NewUsecase(categoryRepository category.CategoryRepository) *Usecase {
	return &Usecase{
		categoryRepository: categoryRepository,
	}
}

func (u *Usecase) run(ctx context.Context, input input) error {
	createInput := category.CreateCategoryInput{
		Name: input.Name,
	}
	if err := createInput.Validate(); err != nil {
		return &handlererror.DisplayableError{
			StatusCode: 400,
			Message:    err.Error(),
			Err:        err,
		}
	}

	categories, err := u.categoryRepository.Search(ctx, category.SearchCategoryCriteria{
		Name: input.Name,
	})
	if err != nil {
		return &handlererror.DisplayableError{
			StatusCode: 500,
			Message:    "カテゴリの作成に失敗しました。",
			Err:        err,
		}
	}
	if len(categories) > 0 {
		return &handlererror.DisplayableError{
			StatusCode: 400,
			Message:    "同じ名前のカテゴリが既に存在しています。",
		}
	}

	if err := u.categoryRepository.Create(ctx, createInput); err != nil {
		return &handlererror.DisplayableError{
			StatusCode: 500,
			Message:    "カテゴリの作成に失敗しました。",
			Err:        err,
		}
	}
	return nil
}
