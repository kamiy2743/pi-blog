package article

import (
	"context"

	domainArticle "blog/internal/domain/article"
)

type ArticleRepositoryStub struct {
	CreateFunc func(ctx context.Context, input domainArticle.CreateArticleInput) (domainArticle.Article, error)
	UpdateFunc func(ctx context.Context, article domainArticle.Article) error
	SearchFunc func(ctx context.Context, criteria domainArticle.SearchArticleCriteria) ([]domainArticle.Article, error)
}

func (s ArticleRepositoryStub) Create(ctx context.Context, input domainArticle.CreateArticleInput) (domainArticle.Article, error) {
	if s.CreateFunc == nil {
		return domainArticle.Article{}, nil
	}
	return s.CreateFunc(ctx, input)
}

func (s ArticleRepositoryStub) Update(ctx context.Context, article domainArticle.Article) error {
	if s.UpdateFunc == nil {
		return nil
	}
	return s.UpdateFunc(ctx, article)
}

func (s ArticleRepositoryStub) Search(ctx context.Context, criteria domainArticle.SearchArticleCriteria) ([]domainArticle.Article, error) {
	if s.SearchFunc == nil {
		return nil, nil
	}
	return s.SearchFunc(ctx, criteria)
}
