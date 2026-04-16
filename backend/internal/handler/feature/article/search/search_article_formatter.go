package search

import "time"

func formatInitial(result initialResult) map[string]any {
	categories := make([]map[string]any, 0, len(result.Categories))
	for _, category := range result.Categories {
		categories = append(categories, map[string]any{
			"id":   category.ID,
			"name": category.Name,
		})
	}

	return map[string]any{
		"categories": categories,
	}
}

func formatPartialSearch(result partialSearchResult) map[string]any {
	articles := make([]map[string]any, 0, len(result.Articles))
	for _, article := range result.Articles {
		categoryNames := make([]string, 0, len(article.Categories))
		for _, category := range article.Categories {
			categoryNames = append(categoryNames, category.Name)
		}

		articles = append(articles, map[string]any{
			"id":            article.ID,
			"title":         article.Title,
			"date":          article.UpdatedAt.Format(time.RFC3339),
			"categoryNames": categoryNames,
		})
	}

	return map[string]any{
		"title":       result.Title,
		"categoryIds": result.CategoryIDs,
		"page":        result.Page,
		"totalCount":  result.TotalCount,
		"totalPages":  result.TotalPages,
		"articles":    articles,
	}
}
