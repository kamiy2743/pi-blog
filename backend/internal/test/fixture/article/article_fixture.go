package article

import (
	"testing"
	"time"

	"blog/internal/ent"
)

type CreateArticleInput struct {
	Title          string
	Body           string
	IsPublished    bool
	PublishStartAt *time.Time
	PublishEndAt   *time.Time
	UpdatedAt      *time.Time
	Categories     []*ent.Category
}

func CreateArticle(
	t *testing.T,
	entClient *ent.Client,
	input CreateArticleInput,
) *ent.Article {
	t.Helper()

	builder := entClient.Article.Create().
		SetTitle(input.Title).
		SetBody(input.Body).
		SetIsPublished(input.IsPublished).
		SetNillablePublishStartAt(input.PublishStartAt).
		SetNillablePublishEndAt(input.PublishEndAt).
		SetNillableUpdatedAt(input.UpdatedAt)

	if len(input.Categories) > 0 {
		builder.AddCategories(input.Categories...)
	}

	model, err := builder.Save(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	return model
}
