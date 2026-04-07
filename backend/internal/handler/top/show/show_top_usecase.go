package show

import (
	"context"

	"blog/internal/domain/article"
	"blog/internal/domain/category"
	"blog/internal/handler/handlererror"
)

type Usecase struct {
	articleRepository  article.ArticleRepository
	categoryRepository category.CategoryRepository
}

func NewUsecase(
	articleRepository article.ArticleRepository,
	categoryRepository category.CategoryRepository,
) *Usecase {
	return &Usecase{
		articleRepository:  articleRepository,
		categoryRepository: categoryRepository,
	}
}

func (u *Usecase) Run(ctx context.Context) (ShowTopResult, *handlererror.DisplayableError) {
	articles, err := u.articleRepository.Search(ctx, article.SearchArticleCriteria{
		Limit:   10,
		OrderBy: article.OrderByLatest,
	})
	if err != nil {
		return ShowTopResult{}, &handlererror.DisplayableError{
			StatusCode:  500,
			Message:     "記事の読み込みに失敗しました。",
			Description: "時間をおいてから、もう一度お試しください。",
			Err:         err,
		}
	}

	categories, err := u.categoryRepository.All(ctx, category.OrderByNameAsc)
	if err != nil {
		return ShowTopResult{}, &handlererror.DisplayableError{
			StatusCode:  500,
			Message:     "カテゴリの読み込みに失敗しました。",
			Description: "時間をおいてから、もう一度お試しください。",
			Err:         err,
		}
	}

	return ShowTopResult{
		LatestArticles: articles,
		Categories:     categories,
	}, nil
}
