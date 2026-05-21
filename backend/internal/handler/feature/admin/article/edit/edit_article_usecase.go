package edit

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

type result struct {
	Categories []category.Category
	Article    article.Article
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

func (u *Usecase) run(ctx context.Context, articleID article.ArticleID) (result, *handlererror.DisplayableError) {
	categories, err := u.categoryRepository.All(ctx, category.OrderByNameAsc)
	if err != nil {
		return result{}, &handlererror.DisplayableError{
			StatusCode: 500,
			Message:    "カテゴリの読み込みに失敗しました。",
			Err:        err,
		}
	}

	articles, err := u.articleRepository.Search(ctx, article.SearchArticleCriteria{
		IDs:                []article.ArticleID{articleID},
		IncludeUnpublished: true,
	})
	if err != nil {
		return result{}, &handlererror.DisplayableError{
			StatusCode: 500,
			Message:    "記事の読み込みに失敗しました。",
			Err:        err,
		}
	}
	if len(articles) == 0 {
		return result{}, &handlererror.DisplayableError{
			StatusCode: 404,
			Message:    "記事が見つかりません。",
		}
	}

	return result{
		Categories: categories,
		Article:    articles[0],
	}, nil
}
