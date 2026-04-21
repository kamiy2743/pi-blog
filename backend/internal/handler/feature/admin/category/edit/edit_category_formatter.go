package edit

import "github.com/romsar/gonertia/v2"

func format(result result) gonertia.Props {
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
