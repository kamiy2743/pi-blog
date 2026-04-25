package show

import (
	"time"

	"github.com/romsar/gonertia/v3"
)

func format(result result) gonertia.Props {
	categoryNames := make([]string, 0, len(result.Article.Categories))
	for _, category := range result.Article.Categories {
		categoryNames = append(categoryNames, category.Name)
	}

	return gonertia.Props{
		"article": gonertia.Props{
			"id":            result.Article.ID,
			"title":         result.Article.Title,
			"body":          result.Article.Body,
			"date":          result.Article.UpdatedAt.Format(time.RFC3339),
			"categoryNames": categoryNames,
		},
	}
}
