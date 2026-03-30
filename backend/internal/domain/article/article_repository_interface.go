package article

import "context"

type ArticleRepository interface {
	Create(ctx context.Context, input CreateArticleInput) (Article, error)
	Update(ctx context.Context, article Article) error
	Search(ctx context.Context, criteria SearchArticleCriteria) ([]Article, error)
}
