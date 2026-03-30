package show

import (
	"context"

	"blog/internal/domain"
	"blog/internal/domain/article"
	"blog/internal/domain/category"
)

func Run(
	ctx context.Context,
	articleRepository article.ArticleRepository,
	categoryRepository category.CategoryRepository,
) (ShowTopResult, error) {
	articles, err := articleRepository.Search(ctx, article.SearchArticleCriteria{
		Limit: 10,
		OrderBy: domain.OrderBy{
			Column:    "updated_at",
			Direction: domain.OrderDirectionDesc,
		},
	})
	if err != nil {
		return ShowTopResult{}, err
	}

	categories, err := categoryRepository.All(ctx)
	if err != nil {
		return ShowTopResult{}, err
	}

	return ShowTopResult{
		LatestArticles: articles,
		Categories:     categories,
	}, nil
}
