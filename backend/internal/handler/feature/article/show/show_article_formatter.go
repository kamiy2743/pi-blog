package show

import (
	"blog/internal/handler/formatter"

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
			"bodyHTML":      result.Article.BodyHTML,
			"date":          formatter.FormatTimeISO8601(&result.Article.UpdatedAt),
			"categoryNames": categoryNames,
		},
	}
}
