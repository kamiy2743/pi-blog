package store

import (
	"context"

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

func (u *Usecase) run(ctx context.Context, input input) error {
	categories, err := u.fetchCategories(ctx, input.CategoryIDs)
	if err != nil {
		return err
	}

	bodyHTML, err := article.RenderMarkdownToHTML(input.Body)
	if err != nil {
		return &handlererror.DisplayableError{
			StatusCode: 400,
			Message:    "本文のHTML変換に失敗しました。",
			Err:        err,
		}
	}

	createInput := article.CreateArticleInput{
		Title:          input.Title,
		BodyMarkdown:   input.Body,
		BodyHTML:       bodyHTML,
		IsPublished:    input.IsPublished,
		PublishStartAt: input.PublishStartAt,
		PublishEndAt:   input.PublishEndAt,
		Categories:     categories,
	}
	if err := createInput.Validate(); err != nil {
		return &handlererror.DisplayableError{
			StatusCode: 400,
			Message:    err.Error(),
			Err:        err,
		}
	}

	if err := u.articleRepository.Create(ctx, createInput); err != nil {
		return &handlererror.DisplayableError{
			StatusCode: 500,
			Message:    "記事の作成に失敗しました。",
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
