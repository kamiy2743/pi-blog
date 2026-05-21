package update

import (
	"context"
	"time"

	"blog/internal/domain/article"
	"blog/internal/domain/category"
	"blog/internal/handler/handlererror"
)

type Usecase struct {
	articleRepository  article.ArticleRepository
	categoryRepository category.CategoryRepository
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

func (u *Usecase) run(ctx context.Context, articleID article.ArticleID, input input) error {
	currentArticles, err := u.articleRepository.Search(ctx, article.SearchArticleCriteria{
		IDs:                []article.ArticleID{articleID},
		IncludeUnpublished: true,
	})
	if err != nil {
		return &handlererror.DisplayableError{
			StatusCode: 500,
			Message:    "記事の更新に失敗しました。",
			Err:        err,
		}
	}
	if len(currentArticles) == 0 {
		return &handlererror.DisplayableError{
			StatusCode: 404,
			Message:    "記事が見つかりません。",
		}
	}

	categories, err := u.fetchCategories(ctx, input.CategoryIDs)
	if err != nil {
		return err
	}

	bodyHTML, err := article.ConvertMarkdownToHTML(input.Body)
	if err != nil {
		return &handlererror.DisplayableError{
			StatusCode: 400,
			Message:    "本文のHTML変換に失敗しました。",
			Err:        err,
		}
	}

	entity := article.Article{
		ID:             articleID,
		Title:          input.Title,
		BodyMarkdown:   input.Body,
		BodyHTML:       bodyHTML,
		IsPublished:    input.IsPublished,
		PublishStartAt: input.PublishStartAt,
		PublishEndAt:   input.PublishEndAt,
		Categories:     categories,
		UpdatedAt:      time.Now(),
	}
	if err := entity.Validate(); err != nil {
		return &handlererror.DisplayableError{
			StatusCode: 400,
			Message:    err.Error(),
			Err:        err,
		}
	}

	if err := u.articleRepository.Update(ctx, entity); err != nil {
		return &handlererror.DisplayableError{
			StatusCode: 500,
			Message:    "記事の更新に失敗しました。",
			Err:        err,
		}
	}
	return nil
}

func (u *Usecase) fetchCategories(ctx context.Context, categoryIDs []category.CategoryID) ([]category.Category, error) {
	if len(categoryIDs) == 0 {
		return []category.Category{}, nil
	}

	categories, err := u.categoryRepository.Search(ctx, category.SearchCategoryCriteria{
		IDs: categoryIDs,
	})
	if err != nil {
		return nil, &handlererror.DisplayableError{
			StatusCode: 500,
			Message:    "カテゴリの読み込みに失敗しました。",
			Err:        err,
		}
	}
	if len(categories) != len(categoryIDs) {
		return nil, &handlererror.DisplayableError{
			StatusCode: 400,
			Message:    "選択したカテゴリが見つかりません。",
		}
	}

	return categories, nil
}
