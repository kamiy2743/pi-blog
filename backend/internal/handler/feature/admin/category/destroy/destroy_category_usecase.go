package destroy

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

func (u *Usecase) run(ctx context.Context, categoryID category.CategoryID) error {
	categories, err := u.categoryRepository.Search(ctx, category.SearchCategoryCriteria{
		IDs: []category.CategoryID{categoryID},
	})
	if err != nil {
		return &handlererror.DisplayableError{
			StatusCode: 500,
			Message:    "カテゴリの削除に失敗しました。",
			Err:        err,
		}
	}
	if len(categories) == 0 {
		return &handlererror.DisplayableError{
			StatusCode: 404,
			Message:    "カテゴリが見つかりません。",
		}
	}

	if err := u.categoryRepository.Delete(ctx, categories[0]); err != nil {
		return &handlererror.DisplayableError{
			StatusCode: 500,
			Message:    "カテゴリの削除に失敗しました。",
			Err:        err,
		}
	}
	return nil
}
