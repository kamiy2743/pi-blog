package article

import (
	"context"
	"time"

	"blog/internal/db/ent"
	entArticle "blog/internal/db/ent/article"
	entCategory "blog/internal/db/ent/category"
	domainArticle "blog/internal/domain/article"
	infraCategory "blog/internal/infra/category"
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
	var applyCategoryOrderErr error
	query := r.client.Article.Query().WithCategories(func(categoryQuery *ent.CategoryQuery) {
		applyCategoryOrderErr = infraCategory.ApplyOrder(categoryQuery, criteria.CategoryOrderBy)
	})
	if applyCategoryOrderErr != nil {
		return nil, applyCategoryOrderErr
	}

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

	var applyCategoryOrderErr error
	articleQuery := r.client.Article.Query().WithCategories(func(categoryQuery *ent.CategoryQuery) {
		applyCategoryOrderErr = infraCategory.ApplyOrder(categoryQuery, criteria.SearchCriteria.CategoryOrderBy)
	})
	if applyCategoryOrderErr != nil {
		return domainArticle.PaginatedArticles{}, applyCategoryOrderErr
	}

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
		for _, categoryID := range criteria.CategoryIDs {
			query.Where(entArticle.HasCategoriesWith(entCategory.ID(uint32(categoryID))))
		}
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

	if err := ApplyOrder(query, criteria.OrderBy); err != nil {
		return err
	}

	return nil
}
