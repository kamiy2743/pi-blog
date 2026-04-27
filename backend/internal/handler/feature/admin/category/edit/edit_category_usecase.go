package edit

import (
	"context"

	"blog/internal/domain/category"
	"blog/internal/handler/handlererror"
)

type Usecase struct {
	categoryRepository category.CategoryRepository
}

type result struct {
	Categories []category.Category
}

func NewUsecase(categoryRepository category.CategoryRepository) *Usecase {
	return &Usecase{
		categoryRepository: categoryRepository,
	}
}

func (u *Usecase) run(ctx context.Context) (result, *handlererror.DisplayableError) {
	categories, err := u.categoryRepository.All(ctx, category.OrderByNameAsc)
	if err != nil {
		return result{}, &handlererror.DisplayableError{
			StatusCode: 500,
			Message:    "カテゴリの読み込みに失敗しました。",
			Err:        err,
		}
	}

	return result{
		Categories: categories,
	}, nil
}
