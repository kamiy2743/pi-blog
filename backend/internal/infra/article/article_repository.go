package article

import (
	"context"
	"fmt"
	"time"

	domainArticle "blog/internal/domain/article"
	"blog/internal/db/ent"
	entArticle "blog/internal/db/ent/article"
	entCategory "blog/internal/db/ent/category"
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

	if err := applySearchCriteria(query, criteria, time.Now()); err != nil {
		return nil, err
	}

	models, err := query.All(ctx)
	if err != nil {
		return nil, err
	}
	return hydrateArticles(models), nil
}

func (r *ArticleRepository) Paginate(ctx context.Context, criteria domainArticle.PaginateArticleCriteria) (domainArticle.PaginatedArticles, error) {
	now := time.Now()

	countQuery := r.client.Article.Query()
	if err := applySearchCriteria(countQuery, criteria.SearchCriteria, now); err != nil {
		return domainArticle.PaginatedArticles{}, err
	}
	totalCount, err := countQuery.Count(ctx)
	if err != nil {
		return domainArticle.PaginatedArticles{}, err
	}

	articleQuery := r.client.Article.Query().WithCategories()

	criteria.SearchCriteria.Limit = &criteria.PerPage
	if err := applySearchCriteria(articleQuery, criteria.SearchCriteria, now); err != nil {
		return domainArticle.PaginatedArticles{}, err
	}
	if criteria.Page > 1 {
		articleQuery.Offset((criteria.Page - 1) * criteria.PerPage)
	}

	models, err := articleQuery.All(ctx)
	if err != nil {
		return domainArticle.PaginatedArticles{}, err
	}

	return domainArticle.PaginatedArticles{
		TotalCount: totalCount,
		Articles:   hydrateArticles(models),
	}, nil
}

func applySearchCriteria(
	query *ent.ArticleQuery,
	criteria domainArticle.SearchArticleCriteria,
	now time.Time,
) error {
	if len(criteria.IDs) > 0 {
		ids := make([]uint32, 0, len(criteria.IDs))
		for _, id := range criteria.IDs {
			ids = append(ids, uint32(id))
		}
		query.Where(entArticle.IDIn(ids...))
	}

	if criteria.Title != "" {
		query.Where(entArticle.TitleContainsFold(criteria.Title))
	}

	if len(criteria.CategoryIDs) > 0 {
		categoryIDs := make([]uint32, 0, len(criteria.CategoryIDs))
		for _, categoryID := range criteria.CategoryIDs {
			categoryIDs = append(categoryIDs, uint32(categoryID))
		}
		query.Where(entArticle.HasCategoriesWith(entCategory.IDIn(categoryIDs...)))
	}

	if !criteria.IncludeUnpublished {
		query.Where(
			entArticle.IsPublished(true),
			entArticle.Or(
				entArticle.PublishStartAtIsNil(),
				entArticle.PublishStartAtLTE(now),
			),
			entArticle.Or(
				entArticle.PublishEndAtIsNil(),
				entArticle.PublishEndAtGTE(now),
			),
		)
	}

	if criteria.Limit != nil {
		query.Limit(*criteria.Limit)
	}

	if criteria.OrderBy != "" {
		switch criteria.OrderBy {
		case domainArticle.OrderByLatest:
			query.Order(
				ent.Desc(entArticle.FieldUpdatedAt),
				ent.Desc(entArticle.FieldID),
			)
		default:
			return fmt.Errorf("未対応の記事の並び順です: %s", criteria.OrderBy)
		}
	}

	return nil
}
