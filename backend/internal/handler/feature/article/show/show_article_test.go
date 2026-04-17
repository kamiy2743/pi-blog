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
	"blog/internal/ent"
	"blog/internal/test"
	fixtureArticle "blog/internal/test/fixture/article"
	fixtureCategory "blog/internal/test/fixture/category"
	"blog/internal/test/helper"
	stubArticle "blog/internal/test/stub/article"
)

func Test記事を表示できる(t *testing.T) {
	initResult := test.Init(t)
	article := setUpRecord(t, initResult.EntClient, true)

	res := callEndpoint(t, initResult.Server, article.ID)

	res.AssertFullProps(t, "article/ShowArticle", map[string]any{
		"article": map[string]any{
			"id":    article.ID,
			"title": "title",
			"body":  "body",
			"date":  "2026-01-01T00:00:00Z",
			"categoryNames": []string{
				"Go",
			},
		},
	})
}

func Test記事IDが不正な場合は404(t *testing.T) {
	initResult := test.Init(t)

	res := helper.RequestInertia(t, initResult.Server, helper.TestInertiaRequest{
		Method: http.MethodGet,
		Path:   "/article/invalid",
	})

	res.AssertNotFound(t)
}

func Test記事が見つからない場合は404(t *testing.T) {
	initResult := test.Init(t)

	res := callEndpoint(t, initResult.Server, 999)

	res.AssertError(t, 404, "記事が見つかりませんでした。", "正しいURLか確認してください。")
}

func Test非公開記事は404(t *testing.T) {
	initResult := test.Init(t)
	article := setUpRecord(t, initResult.EntClient, false)

	res := callEndpoint(t, initResult.Server, article.ID)

	res.AssertError(t, 404, "記事が見つかりませんでした。", "正しいURLか確認してください。")
}

func Test記事の取得に失敗した場合は500(t *testing.T) {
	initResult := test.Init(t, &di.ContainerOptions{
		ArticleRepository: stubArticle.ArticleRepositoryStub{
			SearchFunc: func(ctx context.Context, criteria domainArticle.SearchArticleCriteria) ([]domainArticle.Article, error) {
				return nil, errors.New("test")
			},
		},
	})
	setUpRecord(t, initResult.EntClient, true)

	res := callEndpoint(t, initResult.Server, 1)

	res.AssertError(t, 500, "記事の読み込みに失敗しました。", "時間をおいてから、もう一度お試しください。")
}

func setUpRecord(t *testing.T, entClient *ent.Client, isPublished bool) *ent.Article {
	t.Helper()

	category := fixtureCategory.CreateCategory(t, entClient, fixtureCategory.CreateCategoryInput{Name: "Go"})

	return fixtureArticle.CreateArticle(
		t,
		entClient,
		fixtureArticle.CreateArticleInput{
			Title:       "title",
			Body:        "body",
			IsPublished: isPublished,
			UpdatedAt:   helper.TimePtr(t, "2026-01-01 00:00"),
			Categories: []*ent.Category{
				category,
			},
		},
	)
}

func callEndpoint(t *testing.T, server *httptest.Server, articleID uint32) helper.TestInertiaResponse {
	t.Helper()

	return helper.RequestInertia(t, server, helper.TestInertiaRequest{
		Method: http.MethodGet,
		Path:   fmt.Sprintf("/article/%d", articleID),
	})
}
