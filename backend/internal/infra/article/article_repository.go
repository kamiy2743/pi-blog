package article

import (
	"context"
	"fmt"

	domainArticle "blog/internal/domain/article"
	"blog/internal/ent"
	entArticle "blog/internal/ent/article"
)

type ArticleRepository struct {
	client *ent.Client
}

func NewArticleRepository(client *ent.Client) *ArticleRepository {
	return &ArticleRepository{client: client}
}

func (r *ArticleRepository) Create(ctx context.Context, input domainArticle.CreateArticleInput) (domainArticle.Article, error) {
	return domainArticle.Article{}, nil
}

func (r *ArticleRepository) Update(ctx context.Context, input domainArticle.Article) error {
	return nil
}

func (r *ArticleRepository) Search(ctx context.Context, criteria domainArticle.SearchArticleCriteria) ([]domainArticle.Article, error) {
	query := r.client.Article.Query().WithCategories()

	if criteria.Title != "" {
		query.Where(entArticle.TitleContainsFold(criteria.Title))
	}

	switch criteria.OrderBy {
	case domainArticle.OrderByLatest:
		query.Order(ent.Desc(entArticle.FieldUpdatedAt))
	default:
		return nil, fmt.Errorf("未対応の記事の並び順です: %s", criteria.OrderBy)
	}

	if criteria.Limit > 0 {
		query.Limit(criteria.Limit)
	}

	models, err := query.All(ctx)
	if err != nil {
		return nil, err
	}
	return hydrateArticles(models), nil
}
