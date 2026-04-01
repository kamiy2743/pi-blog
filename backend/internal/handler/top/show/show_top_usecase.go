package show

import (
	"context"

	"blog/internal/domain/article"
	"blog/internal/domain/category"
)

func Run(
	ctx context.Context,
	articleRepository article.ArticleRepository,
	categoryRepository category.CategoryRepository,
) (ShowTopResult, error) {
	articles, err := articleRepository.Search(ctx, article.SearchArticleCriteria{
		Limit:   10,
		OrderBy: article.OrderByLatest,
	})
	if err != nil {
		return ShowTopResult{}, err
	}

	categories, err := categoryRepository.All(ctx, category.OrderByNameAsc)
	if err != nil {
		return ShowTopResult{}, err
	}

	return ShowTopResult{
		LatestArticles: articles,
		Categories:     categories,
	}, nil
}
