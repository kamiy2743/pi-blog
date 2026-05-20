package article

import (
	"errors"
	"strings"
	"time"

	"blog/internal/domain/category"
)

type Article struct {
	ID             ArticleID
	Title          string
	BodyMarkdown   string
	BodyHTML       string
	IsPublished    bool
	PublishStartAt *time.Time
	PublishEndAt   *time.Time
	Categories     []category.Category
	UpdatedAt      time.Time
}

type CreateArticleInput struct {
	Title          string
	BodyMarkdown   string
	BodyHTML       string
	IsPublished    bool
	PublishStartAt *time.Time
	PublishEndAt   *time.Time
	Categories     []category.Category
}

func (a Article) Validate() error {
	if err := validateContent(a.Title, a.BodyMarkdown, a.BodyHTML, a.PublishStartAt, a.PublishEndAt, a.Categories); err != nil {
		return err
	}
	if a.UpdatedAt.IsZero() {
		return errors.New("更新日時が不正です")
	}
	return nil
}

func (a CreateArticleInput) Validate() error {
	if err := validateContent(a.Title, a.BodyMarkdown, a.BodyHTML, a.PublishStartAt, a.PublishEndAt, a.Categories); err != nil {
		return err
	}
	return nil
}

func validateContent(
	title string,
	bodyMarkdown string,
	bodyHTML string,
	publishStartAt *time.Time,
	publishEndAt *time.Time,
	categories []category.Category,
) error {
	if strings.TrimSpace(title) == "" {
		return errors.New("記事タイトルは必須です")
	}
	if strings.TrimSpace(bodyMarkdown) == "" {
		return errors.New("記事本文は必須です")
	}
	if strings.TrimSpace(bodyHTML) == "" {
		return errors.New("記事本文HTMLは必須です")
	}
	if publishStartAt != nil && publishEndAt != nil && publishEndAt.Before(*publishStartAt) {
		return errors.New("公開終了時刻は公開開始時刻以降を指定してください")
	}
	for _, category := range categories {
		if err := category.Validate(); err != nil {
			return err
		}
	}
	return nil
}
