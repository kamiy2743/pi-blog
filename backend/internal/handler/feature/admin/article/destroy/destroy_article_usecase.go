package destroy

import (
	"context"

	"blog/internal/domain/article"
	"blog/internal/handler/handlererror"
)

type Usecase struct {
	articleRepository article.ArticleRepository
}

func NewUsecase(articleRepository article.ArticleRepository) *Usecase {
	return &Usecase{
		articleRepository: articleRepository,
	}
}

func (u *Usecase) run(ctx context.Context, articleID article.ArticleID) error {
	articles, err := u.articleRepository.Search(ctx, article.SearchArticleCriteria{
		IDs:                []article.ArticleID{articleID},
		IncludeUnpublished: true,
	})
	if err != nil {
		return &handlererror.DisplayableError{
			StatusCode: 500,
			Message:    "記事の削除に失敗しました。",
			Err:        err,
		}
	}
	if len(articles) == 0 {
		return &handlererror.DisplayableError{
			StatusCode: 404,
			Message:    "記事が見つかりません。",
		}
	}

	if err := u.articleRepository.Delete(ctx, articles[0]); err != nil {
		return &handlererror.DisplayableError{
			StatusCode: 500,
			Message:    "記事の削除に失敗しました。",
			Err:        err,
		}
	}
	return nil
}
