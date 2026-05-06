package update

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

func (u *Usecase) run(ctx context.Context, categoryID category.CategoryID, input input) error {
	entity := category.Category{
		ID:   categoryID,
		Name: input.Name,
	}
	if err := entity.Validate(); err != nil {
		return &handlererror.DisplayableError{
			StatusCode: 400,
			Message:    err.Error(),
			Err:        err,
		}
	}

	currentCategories, err := u.categoryRepository.Search(ctx, category.SearchCategoryCriteria{
		IDs: []category.CategoryID{categoryID},
	})
	if err != nil {
		return &handlererror.DisplayableError{
			StatusCode: 500,
			Message:    "カテゴリの更新に失敗しました。",
			Err:        err,
		}
	}
	if len(currentCategories) == 0 {
		return &handlererror.DisplayableError{
			StatusCode: 404,
			Message:    "カテゴリが見つかりません。",
		}
	}

	sameNameCategories, err := u.categoryRepository.Search(ctx, category.SearchCategoryCriteria{
		Name: input.Name,
	})
	if err != nil {
		return &handlererror.DisplayableError{
			StatusCode: 500,
			Message:    "カテゴリの更新に失敗しました。",
			Err:        err,
		}
	}
	if len(sameNameCategories) > 0 {
		return &handlererror.DisplayableError{
			StatusCode: 400,
			Message:    "同じ名前のカテゴリが既に存在しています。",
		}
	}

	if err := u.categoryRepository.Update(ctx, entity); err != nil {
		return &handlererror.DisplayableError{
			StatusCode: 500,
			Message:    "カテゴリの更新に失敗しました。",
			Err:        err,
		}
	}
	return nil
}
