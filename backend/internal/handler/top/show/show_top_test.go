package show_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"blog/internal/di"
	domainArticle "blog/internal/domain/article"
	domainCategory "blog/internal/domain/category"
	"blog/internal/ent"
	"blog/internal/test"
	fixtureArticle "blog/internal/test/fixture/article"
	fixtureCategory "blog/internal/test/fixture/category"
	"blog/internal/test/helper"
	stubArticle "blog/internal/test/stub/article"
	stubCategory "blog/internal/test/stub/category"
)

func Testトップページを表示できる(t *testing.T) {
	initResult := test.Init(t)
	setUpRecords(t, initResult.EntClient)

	res := callEndpoint(t, initResult.Server)
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

func Test記事の取得に失敗した場合は500(t *testing.T) {
	initResult := test.Init(t, &di.ContainerOptions{
		ArticleRepository: stubArticle.ArticleRepositoryStub{
			SearchFunc: func(ctx context.Context, criteria domainArticle.SearchArticleCriteria) ([]domainArticle.Article, error) {
				return nil, errors.New("test")
			},
		},
	})
	setUpRecords(t, initResult.EntClient)

	res := callEndpoint(t, initResult.Server)
	res.AssertError(t, http.StatusInternalServerError, "記事の読み込みに失敗しました。", "時間をおいてから、もう一度お試しください。")
}

func Testカテゴリの取得に失敗した場合は500(t *testing.T) {
	initResult := test.Init(t, &di.ContainerOptions{
		CategoryRepository: stubCategory.CategoryRepositoryStub{
			AllFunc: func(ctx context.Context, orderBy domainCategory.OrderBy) ([]domainCategory.Category, error) {
				return nil, errors.New("test")
			},
		},
	})
	setUpRecords(t, initResult.EntClient)

	res := callEndpoint(t, initResult.Server)
	res.AssertError(t, http.StatusInternalServerError, "カテゴリの読み込みに失敗しました。", "時間をおいてから、もう一度お試しください。")
}

func setUpRecords(t *testing.T, entClient *ent.Client) {
	t.Helper()

	category1 := fixtureCategory.CreateCategory(t, entClient, fixtureCategory.CreateCategoryInput{Name: "category-b"})
	category2 := fixtureCategory.CreateCategory(t, entClient, fixtureCategory.CreateCategoryInput{Name: "category-a"})

	for i := 0; i < 2; i++ {
		fixtureArticle.CreateArticle(
			t,
			entClient,
			fixtureArticle.CreateArticleInput{
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

func callEndpoint(t *testing.T, server *httptest.Server) helper.TestInertiaResponse {
	return helper.RequestInertia(t, server, helper.TestInertiaRequest{
		Method: http.MethodGet,
		Path:   "/",
	})
}
