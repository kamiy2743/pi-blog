package article

import (
	"blog/internal/db/ent"
	domainArticle "blog/internal/domain/article"
	infraCategory "blog/internal/infra/category"
)

func hydrateArticle(model *ent.Article) domainArticle.Article {
	return domainArticle.Article{
		ID:             domainArticle.ArticleID(model.ID),
		Title:          model.Title,
		Body:           model.Body,
		IsPublished:    model.IsPublished,
		PublishStartAt: model.PublishStartAt,
		PublishEndAt:   model.PublishEndAt,
		Categories:     infraCategory.HydrateCategories(model.Edges.Categories),
		UpdatedAt:      model.UpdatedAt,
	}
}

func hydrateArticles(models []*ent.Article) []domainArticle.Article {
	articles := make([]domainArticle.Article, 0, len(models))
	for _, model := range models {
		articles = append(articles, hydrateArticle(model))
	}
	return articles
}
