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

type records struct {
	Articles   []*ent.Article
	Categories []*ent.Category
}

func Testトップページを表示できる(t *testing.T) {
	initResult := test.Init(t)
	records := setUpRecords(t, initResult.EntClient)

	res := callEndpoint(t, initResult.Server)

	res.AssertProps(t, "ShowTop", map[string]any{
		"latestArticles": []map[string]any{
			{
				"id":    records.Articles[1].ID,
				"title": "title1",
				"date":  "2026-01-02T00:00:00Z",
				"categoryNames": []string{
					"category-b",
					"category-a",
				},
			},
			{
				"id":    records.Articles[0].ID,
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
				"id":   records.Categories[0].ID,
				"name": "category-a",
			},
			{
				"id":   records.Categories[1].ID,
				"name": "category-b",
			},
		},
	})
}

func Test最大10件までしか表示されない(t *testing.T) {
	initResult := test.Init(t)
	setUpRecords(t, initResult.EntClient)

	for i := 0; i < 10; i++ {
		fixtureArticle.CreateArticle(
			t,
			initResult.EntClient,
			fixtureArticle.CreateArticleInput{
				IsPublished: true,
			},
		)
	}

	res := callEndpoint(t, initResult.Server)

	res.AssertPropsCount(t, "ShowTop", "latestArticles", 10)
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

func setUpRecords(t *testing.T, entClient *ent.Client) records {
	t.Helper()

	categoryB := fixtureCategory.CreateCategory(t, entClient, fixtureCategory.CreateCategoryInput{Name: "category-b"})
	categoryA := fixtureCategory.CreateCategory(t, entClient, fixtureCategory.CreateCategoryInput{Name: "category-a"})

	articles := make([]*ent.Article, 0, 2)
	for i := 0; i < 2; i++ {
		articles = append(articles, fixtureArticle.CreateArticle(
			t,
			entClient,
			fixtureArticle.CreateArticleInput{
				Title:       fmt.Sprintf("title%d", i),
				Body:        "content",
				IsPublished: true,
				UpdatedAt:   helper.TimePtr(t, fmt.Sprintf("2026-01-0%d 00:00", i+1)),
				Categories: []*ent.Category{
					categoryB,
					categoryA,
				},
			},
		))
	}
	articles = append(articles, fixtureArticle.CreateArticle(
		t,
		entClient,
		fixtureArticle.CreateArticleInput{
			Title:       "unpublished",
			IsPublished: false,
		},
	))

	return records{
		Articles: articles,
		Categories: []*ent.Category{
			categoryA,
			categoryB,
		},
	}
}

func callEndpoint(t *testing.T, server *httptest.Server) helper.TestInertiaResponse {
	return helper.RequestInertia(t, server, helper.TestInertiaRequest{
		Method: http.MethodGet,
		Path:   "/",
	})
}
