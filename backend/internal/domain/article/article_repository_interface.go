package article

import "context"

type ArticleRepository interface {
	Create(ctx context.Context, input CreateArticleInput) error
	Update(ctx context.Context, article Article) error
	Search(ctx context.Context, criteria SearchArticleCriteria) ([]Article, error)
	Paginate(ctx context.Context, criteria PaginateArticleCriteria) (PaginatedArticles, error)
}
