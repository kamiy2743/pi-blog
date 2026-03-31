package show

import (
	"time"

	"github.com/romsar/gonertia/v2"
)

func Format(result ShowTopResult) gonertia.Props {
	latestArticles := make([]map[string]any, 0, len(result.LatestArticles))
	for _, article := range result.LatestArticles {
		categoryNames := make([]string, 0, len(article.Categories))
		for _, category := range article.Categories {
			categoryNames = append(categoryNames, category.Name)
		}
		latestArticles = append(latestArticles, map[string]any{
			"id":            article.ID,
			"title":         article.Title,
			"date":          article.UpdatedAt.Format(time.RFC3339),
			"categoryNames": categoryNames,
		})
	}

	categories := make([]map[string]any, 0, len(result.Categories))
	for _, category := range result.Categories {
		categories = append(categories, map[string]any{
			"id":   category.ID,
			"name": category.Name,
		})
	}

	return gonertia.Props{
		"latestArticles": latestArticles,
		"categories":     categories,
	}
}
