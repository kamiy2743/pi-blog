package search

import (
	"context"
	"math"

	"blog/internal/domain/article"
	"blog/internal/domain/category"
	"blog/internal/handler/handlererror"
)

type Usecase struct {
	articleRepository  article.ArticleRepository
	categoryRepository category.CategoryRepository
}

type input struct {
	Title       string
	CategoryIDs []category.CategoryID
	Page        int
	PerPage     int
}

type initialResult struct {
	Categories []category.Category
}

type partialSearchResult struct {
	Title       string
	CategoryIDs []category.CategoryID
	Page        int
	TotalCount  int
	TotalPages  int
	Articles    []article.Article
}

func NewUsecase(
	articleRepository article.ArticleRepository,
	categoryRepository category.CategoryRepository,
) *Usecase {
	return &Usecase{
		articleRepository:  articleRepository,
		categoryRepository: categoryRepository,
	}
}

func (u *Usecase) runInitial(ctx context.Context) (initialResult, *handlererror.DisplayableError) {
	categories, err := u.categoryRepository.All(ctx, category.OrderByNameAsc)
	if err != nil {
		return initialResult{}, &handlererror.DisplayableError{
			StatusCode: 500,
			Message:    "カテゴリの読み込みに失敗しました。",
			Err:        err,
		}
	}

	return initialResult{
		Categories: categories,
	}, nil
}

func (u *Usecase) runPartialSearch(ctx context.Context, input input) (partialSearchResult, *handlererror.DisplayableError) {
	validCategoryIDs := []category.CategoryID{}
	if len(input.CategoryIDs) > 0 {
		categories, err := u.categoryRepository.Search(ctx, category.SearchCategoryCriteria{
			IDs:     input.CategoryIDs,
			OrderBy: category.OrderByNameAsc,
		})
		if err != nil {
			return partialSearchResult{}, &handlererror.DisplayableError{
				StatusCode: 500,
				Message:    "カテゴリの読み込みに失敗しました。",
				Err:        err,
			}
		}
		for _, category := range categories {
			validCategoryIDs = append(validCategoryIDs, category.ID)
		}
	}

	paginatedArticles, err := u.articleRepository.Paginate(ctx, article.PaginateArticleCriteria{
		SearchCriteria: article.SearchArticleCriteria{
			Title:       input.Title,
			CategoryIDs: validCategoryIDs,
			OrderBy:     article.OrderByLatest,
		},
		Page:    input.Page,
		PerPage: input.PerPage,
	})
	if err != nil {
		return partialSearchResult{}, &handlererror.DisplayableError{
			StatusCode: 500,
			Message:    "記事の読み込みに失敗しました。",
			Err:        err,
		}
	}

	return partialSearchResult{
		Title:       input.Title,
		CategoryIDs: validCategoryIDs,
		Page:        input.Page,
		TotalCount:  paginatedArticles.TotalCount,
		TotalPages:  int(math.Ceil(float64(paginatedArticles.TotalCount) / float64(input.PerPage))),
		Articles:    paginatedArticles.Articles,
	}, nil
}
