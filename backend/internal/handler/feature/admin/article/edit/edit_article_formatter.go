package edit

import (
	"blog/internal/datetime"

	"github.com/romsar/gonertia/v3"
)

func format(result result) gonertia.Props {
	categoryIDs := make([]uint32, 0, len(result.Article.Categories))
	for _, category := range result.Article.Categories {
		categoryIDs = append(categoryIDs, uint32(category.ID))
	}

	categories := make([]gonertia.Props, 0, len(result.Categories))
	for _, category := range result.Categories {
		categories = append(categories, gonertia.Props{
			"id":   category.ID,
			"name": category.Name,
		})
	}

	return gonertia.Props{
		"categories": categories,
		"article": gonertia.Props{
			"id":             result.Article.ID,
			"title":          result.Article.Title,
			"bodyMarkdown":   result.Article.BodyMarkdown,
			"isPublished":    result.Article.IsPublished,
			"publishStartAt": datetime.FormatISO8601(result.Article.PublishStartAt),
			"publishEndAt":   datetime.FormatISO8601(result.Article.PublishEndAt),
			"categoryIds":    categoryIDs,
		},
	}
}
