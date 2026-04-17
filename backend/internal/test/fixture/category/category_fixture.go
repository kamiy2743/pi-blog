package category

import (
	"testing"

	"blog/internal/db/ent"
)

type CreateCategoryInput struct {
	Name string
}

func CreateCategory(
	t *testing.T,
	entClient *ent.Client,
	input CreateCategoryInput,
) *ent.Category {
	t.Helper()

	model, err := entClient.Category.Create().
		SetName(input.Name).
		Save(t.Context())

	if err != nil {
		t.Fatal(err)
	}
	return model
}
