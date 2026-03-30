package article

import (
	"context"
	"slices"
	"strings"
	"time"

	"blog/internal/domain"
	"blog/internal/domain/article"
	"blog/internal/domain/category"
)

type ArticleRepository struct{}

func NewArticleRepository() *ArticleRepository {
	return &ArticleRepository{}
}

func (r *ArticleRepository) Create(_ context.Context, input article.CreateArticleInput) (article.Article, error) {
	if err := input.Validate(); err != nil {
		return article.Article{}, err
	}

	now := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	raspiCategoryID, _ := category.ParseCategoryID("1")
	return article.Article{
		ID:         article.ArticleID("1"),
		Title:      input.Title,
		ContentMD:  input.ContentMD,
		Categories: withDefaultCategories(input.Categories, []category.Category{{ID: raspiCategoryID, Name: "Raspberry Pi"}}),
		UpdatedAt:  now,
	}, nil
}

func (r *ArticleRepository) Search(_ context.Context, criteria article.SearchArticleCriteria) ([]article.Article, error) {
	raspiCategoryID, _ := category.ParseCategoryID("1")
	infraCategoryID, _ := category.ParseCategoryID("2")
	articles := []article.Article{
		{
			ID:        article.ArticleID("1"),
			Title:     "Raspberry Pi 5 の常時稼働メモ",
			ContentMD: "# Raspberry Pi 5\n\n常時稼働の設定メモです。",
			Categories: []category.Category{
				{ID: raspiCategoryID, Name: "Raspberry Pi"},
			},
			UpdatedAt: time.Date(2025, 1, 2, 12, 0, 0, 0, time.UTC),
		},
		{
			ID:        article.ArticleID("2"),
			Title:     "Cloudflare Tunnel で公開構成を整理する",
			ContentMD: "# Cloudflare Tunnel\n\n公開構成の整理メモです。",
			Categories: []category.Category{
				{ID: infraCategoryID, Name: "Infrastructure"},
			},
			UpdatedAt: time.Date(2025, 1, 4, 12, 0, 0, 0, time.UTC),
		},
	}

	if strings.TrimSpace(criteria.Title) == "" {
		return applySearchCriteria(articles, criteria), nil
	}

	filtered := make([]article.Article, 0, len(articles))
	for _, article := range articles {
		if strings.Contains(article.Title, criteria.Title) {
			filtered = append(filtered, article)
		}
	}
	return applySearchCriteria(filtered, criteria), nil
}

func (r *ArticleRepository) Update(_ context.Context, article article.Article) error {
	return article.Validate()
}

func withDefaultCategories(categories []category.Category, defaultCategories []category.Category) []category.Category {
	if len(categories) == 0 {
		return defaultCategories
	}
	return categories
}

func applySearchCriteria(articles []article.Article, criteria article.SearchArticleCriteria) []article.Article {
	filtered := slices.Clone(articles)

	if criteria.OrderBy.Column != "" {
		sortArticles(filtered, criteria.OrderBy)
	}

	if criteria.Limit > 0 && len(filtered) > criteria.Limit {
		filtered = filtered[:criteria.Limit]
	}

	return filtered
}

func sortArticles(articles []article.Article, orderBy domain.OrderBy) {
	switch orderBy.Column {
	case "updated_at":
		slices.SortFunc(articles, func(a article.Article, b article.Article) int {
			if orderBy.Direction == domain.OrderDirectionAsc {
				return a.UpdatedAt.Compare(b.UpdatedAt)
			}
			return b.UpdatedAt.Compare(a.UpdatedAt)
		})
	}
}
