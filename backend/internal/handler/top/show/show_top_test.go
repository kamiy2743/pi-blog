package show_test

import (
	"fmt"
	"net/http"
	"testing"

	"blog/internal/ent"
	"blog/internal/test"
	"blog/internal/test/fixture/article"
	"blog/internal/test/fixture/category"
	"blog/internal/test/helper"
)

func Testトップページを表示できる(t *testing.T) {
	initResult := test.Init(t)
	setUpRecords(t, initResult.EntClient)

	res := helper.RequestInertia(t, initResult.Server, helper.TestInertiaRequest{
		Method: http.MethodGet,
		Path:   "/",
	})

	res.AssertProps(t, "ShowTop", map[string]any{
		"errors": map[string]any{},
		"latestArticles": []map[string]any{
			{
				"id":    2,
				"title": "title1",
				"date":  "2026-01-02T00:00:00Z",
				"categoryNames": []string{
					"category-b",
					"category-a",
				},
			},
			{
				"id":    1,
				"title": "title0",
				"date":  "2026-01-01T00:00:00Z",
				"categoryNames": []string{
					"category-b",
					"category-a",
				},
			},
		},
		"categories": []map[string]any{
			{
				"id":   2,
				"name": "category-a",
			},
			{
				"id":   1,
				"name": "category-b",
			},
		},
	})
}

func setUpRecords(t *testing.T, entClient *ent.Client) {
	t.Helper()

	category1 := category.CreateCategory(t, entClient, category.CreateCategoryInput{Name: "category-b"})
	category2 := category.CreateCategory(t, entClient, category.CreateCategoryInput{Name: "category-a"})

	for i := 0; i < 2; i++ {
		article.CreateArticle(
			t,
			entClient,
			article.CreateArticleInput{
				Title:       fmt.Sprintf("title%d", i),
				Body:        "content",
				IsPublished: true,
				UpdatedAt:   helper.TimePtr(t, fmt.Sprintf("2026-01-0%d 00:00", i+1)),
				Categories: []*ent.Category{
					category1,
					category2,
				},
			},
		)
	}
}
