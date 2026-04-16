package article

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"blog/internal/domain/category"
)

type Article struct {
	ID             ArticleID
	Title          string
	Body           string
	IsPublished    bool
	PublishStartAt *time.Time
	PublishEndAt   *time.Time
	Categories     []category.Category
	UpdatedAt      time.Time
}

type CreateArticleInput struct {
	Title      string
	Body       string
	Categories []category.Category
}

var (
	errInvalidArticle   = errors.New("記事が不正です")
	errEmptyTitle       = errors.New("記事タイトルは必須です")
	errEmptyBody        = errors.New("記事本文は必須です")
	errInvalidUpdatedAt = errors.New("更新日時が不正です")
)

func (a Article) Validate() error {
	if err := validateContent(a.Title, a.Body, a.Categories); err != nil {
		return err
	}
	if a.UpdatedAt.IsZero() {
		return fmt.Errorf("%w: %w", errInvalidArticle, errInvalidUpdatedAt)
	}
	return nil
}

func (a CreateArticleInput) Validate() error {
	if err := validateContent(a.Title, a.Body, a.Categories); err != nil {
		return err
	}
	return nil
}

func validateContent(title string, body string, categories []category.Category) error {
	if strings.TrimSpace(title) == "" {
		return fmt.Errorf("%w: %w", errInvalidArticle, errEmptyTitle)
	}
	if strings.TrimSpace(body) == "" {
		return fmt.Errorf("%w: %w", errInvalidArticle, errEmptyBody)
	}
	for _, category := range categories {
		if err := category.Validate(); err != nil {
			return fmt.Errorf("%w: %w", errInvalidArticle, err)
		}
	}
	return nil
}
