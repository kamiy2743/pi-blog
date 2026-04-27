package show

import (
	"context"

	"blog/internal/domain/article"
	"blog/internal/handler/handlererror"
)

type Usecase struct {
	articleRepository article.ArticleRepository
}

type result struct {
	Article article.Article
}

func NewUsecase(articleRepository article.ArticleRepository) *Usecase {
	return &Usecase{
		articleRepository: articleRepository,
	}
}

func (u *Usecase) run(ctx context.Context, articleID article.ArticleID) (result, *handlererror.DisplayableError) {
	articles, err := u.articleRepository.Search(ctx, article.SearchArticleCriteria{
		IDs: []article.ArticleID{articleID},
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
			Message:    "記事が見つかりませんでした。",
		}
	}

	return result{
		Article: articles[0],
	}, nil
}
