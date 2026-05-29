package comment

import (
	"testing"

	"blog/internal/db/ent"
)

type CreateCommentInput struct {
	Article    *ent.Article
	AuthorName string
	Body       string
	IsVisible  *bool
}

func CreateComment(
	t *testing.T,
	entClient *ent.Client,
	input CreateCommentInput,
) *ent.Comment {
	t.Helper()

	builder := entClient.Comment.Create().
		SetArticle(input.Article).
		SetAuthorName(input.AuthorName).
		SetBody(input.Body).
		SetNillableIsVisible(input.IsVisible)

	model, err := builder.Save(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	return model
}
