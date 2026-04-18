package edit

import "github.com/romsar/gonertia/v2"

func format(result result) gonertia.Props {
	categories := make([]map[string]any, 0, len(result.Categories))
	for _, category := range result.Categories {
		categories = append(categories, map[string]any{
			"id":   category.ID,
			"name": category.Name,
		})
	}

	return gonertia.Props{
		"categories": categories,
	}
}
