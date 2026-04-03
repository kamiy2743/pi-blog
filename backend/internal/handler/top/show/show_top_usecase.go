package show

import (
	"context"

	"blog/internal/domain/article"
	"blog/internal/domain/category"
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

func (u *Usecase) Run(ctx context.Context) (ShowTopResult, error) {
	articles, err := u.articleRepository.Search(ctx, article.SearchArticleCriteria{
		Limit:   10,
		OrderBy: article.OrderByLatest,
	})
	if err != nil {
		return ShowTopResult{}, err
	}

	categories, err := u.categoryRepository.All(ctx, category.OrderByNameAsc)
	if err != nil {
		return ShowTopResult{}, err
	}

	return ShowTopResult{
		LatestArticles: articles,
		Categories:     categories,
	}, nil
}
