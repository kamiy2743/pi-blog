package search

import (
	"time"

	"github.com/romsar/gonertia/v3"
)

func formatInitial(result initialResult) gonertia.Props {
	categories := make([]gonertia.Props, 0, len(result.Categories))
	for _, category := range result.Categories {
		categories = append(categories, gonertia.Props{
			"id":   category.ID,
			"name": category.Name,
		})
	}

	return gonertia.Props{
		"categories": categories,
	}
}

func formatPartialSearch(result partialSearchResult) gonertia.Props {
	articles := make([]gonertia.Props, 0, len(result.Articles))
	for _, article := range result.Articles {
		categoryNames := make([]string, 0, len(article.Categories))
		for _, category := range article.Categories {
			categoryNames = append(categoryNames, category.Name)
		}

		articles = append(articles, gonertia.Props{
			"id":            article.ID,
			"title":         article.Title,
			"date":          article.UpdatedAt.Format(time.RFC3339),
			"categoryNames": categoryNames,
		})
	}

	return gonertia.Props{
		"title":       result.Title,
		"categoryIds": result.CategoryIDs,
		"page":        result.Page,
		"totalCount":  result.TotalCount,
		"totalPages":  result.TotalPages,
		"articles":    articles,
	}
}
