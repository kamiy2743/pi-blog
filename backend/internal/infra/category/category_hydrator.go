package category

import (
	"blog/internal/db/ent"
	domainCategory "blog/internal/domain/category"
)

func hydrateCategory(model *ent.Category) domainCategory.Category {
	return domainCategory.Category{
		ID:   domainCategory.CategoryID(model.ID),
		Name: model.Name,
	}
}

func HydrateCategories(models []*ent.Category) []domainCategory.Category {
	categories := make([]domainCategory.Category, 0, len(models))
	for _, model := range models {
		categories = append(categories, hydrateCategory(model))
	}
	return categories
}
