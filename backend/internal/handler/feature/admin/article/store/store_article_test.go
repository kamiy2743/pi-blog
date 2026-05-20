package store_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"blog/internal/db/ent"
	entArticle "blog/internal/db/ent/article"
	"blog/internal/di"
	domainArticle "blog/internal/domain/article"
	domainCategory "blog/internal/domain/category"
	"blog/internal/test"
	fixtureCategory "blog/internal/test/fixture/category"
	"blog/internal/test/helper"
	inertiaAction "blog/internal/test/helper/inertia/action"
	stubArticle "blog/internal/test/stub/article"
	stubCategory "blog/internal/test/stub/category"
)

func Test記事を作成できる(t *testing.T) {
	initResult := test.Init(t)
	categoryDocker := fixtureCategory.CreateCategory(t, initResult.EntClient, fixtureCategory.CreateCategoryInput{Name: "Docker"})
	categoryGo := fixtureCategory.CreateCategory(t, initResult.EntClient, fixtureCategory.CreateCategoryInput{Name: "Go"})

	res := callEndpoint(t, initResult.Server, map[string]any{
		"title":          "Go on Raspberry Pi",
		"body":           "# hello",
		"isPublished":    "true",
		"publishStartAt": "2026-01-02T03:04",
		"publishEndAt":   "2026-01-03T04:05",
		"categoryIds":    []string{categoryIDString(categoryDocker.ID), categoryIDString(categoryGo.ID)},
	})

	res.AssertRedirectTo(t, "/admin")

	assertArticleCount(t, initResult.EntClient, 1)

	articles := fetchArticles(t, initResult.EntClient)
	helper.AssertEqual(t, "Go on Raspberry Pi", articles[0].Title, "タイトルが不正です")
	helper.AssertEqual(t, "# hello", articles[0].Body, "本文が不正です")
	helper.AssertEqual(t, true, articles[0].IsPublished, "公開状態が不正です")
	helper.AssertEqual(t, helper.TimePtr(t, "2026-01-02 03:04"), articles[0].PublishStartAt, "公開開始時刻が不正です")
	helper.AssertEqual(t, helper.TimePtr(t, "2026-01-03 04:05"), articles[0].PublishEndAt, "公開終了時刻が不正です")
	helper.AssertEqual(t, []uint32{categoryDocker.ID, categoryGo.ID}, getArticleCategoryIDs(articles[0]), "カテゴリIDが不正です")
}

func Testタイトルが空ならバリデーションエラー(t *testing.T) {
	initResult := test.Init(t)

	res := callEndpoint(t, initResult.Server, map[string]any{
		"title":          "",
		"body":           "content",
		"isPublished":    "false",
		"publishStartAt": "",
		"publishEndAt":   "",
		"categoryIds":    []string{"1", "2"},
	})

	res.AssertRedirectTo(t, "/admin/article/new")
	res.AssertOldInput(t, initResult.Server, initResult.SessionManager, map[string]string{
		"title":          "",
		"body":           "content",
		"isPublished":    "false",
		"publishStartAt": "",
		"publishEndAt":   "",
		"categoryIds":    "1,2",
	})
	res.AssertValidationError(t, initResult.Server, initResult.SessionManager, handlererror.ValidationErrorMessages{
		"title": "タイトルを入力してください。",
	})

	assertArticleCount(t, initResult.EntClient, 0)
}

func Test本文が空ならバリデーションエラー(t *testing.T) {
	initResult := test.Init(t)

	res := callEndpoint(t, initResult.Server, map[string]any{
		"title":          "title",
		"body":           "",
		"isPublished":    "false",
		"publishStartAt": "",
		"publishEndAt":   "",
		"categoryIds":    []string{},
	})

	res.AssertRedirectTo(t, "/admin/article/new")
	res.AssertOldInput(t, initResult.Server, initResult.SessionManager, map[string]string{
		"title":          "title",
		"body":           "",
		"isPublished":    "false",
		"publishStartAt": "",
		"publishEndAt":   "",
		"categoryIds":    "",
	})
	res.AssertValidationError(t, initResult.Server, initResult.SessionManager, handlererror.ValidationErrorMessages{
		"body": "本文を入力してください。",
	})

	assertArticleCount(t, initResult.EntClient, 0)
}

func Test公開状態がboolでなければバリデーションエラー(t *testing.T) {
	initResult := test.Init(t)

	res := callEndpoint(t, initResult.Server, map[string]any{
		"title":          "title",
		"body":           "content",
		"isPublished":    "invalid",
		"publishStartAt": "",
		"publishEndAt":   "",
		"categoryIds":    []string{},
	})

	res.AssertRedirectTo(t, "/admin/article/new")
	res.AssertValidationError(t, initResult.Server, initResult.SessionManager, handlererror.ValidationErrorMessages{
		"isPublished": "公開状態が不正です。",
	})

	assertArticleCount(t, initResult.EntClient, 0)
}

func Test公開終了時刻が公開開始時刻より前ならバリデーションエラー(t *testing.T) {
	initResult := test.Init(t)

	res := callEndpoint(t, initResult.Server, map[string]any{
		"title":          "title",
		"body":           "content",
		"isPublished":    "false",
		"publishStartAt": "2026-01-03T04:05",
		"publishEndAt":   "2026-01-02T03:04",
		"categoryIds":    []string{},
	})

	res.AssertRedirectTo(t, "/admin/article/new")
	res.AssertOldInput(t, initResult.Server, initResult.SessionManager, map[string]string{
		"title":          "title",
		"body":           "content",
		"isPublished":    "false",
		"publishStartAt": "2026-01-03T04:05",
		"publishEndAt":   "2026-01-02T03:04",
		"categoryIds":    "",
	})
	res.AssertValidationError(t, initResult.Server, initResult.SessionManager, handlererror.ValidationErrorMessages{
		"publishStartAt": "公開開始時刻は公開終了時刻より前を指定してください。",
	})

	assertArticleCount(t, initResult.EntClient, 0)
}

func Test複数のバリデーションエラーをまとめて返す(t *testing.T) {
	initResult := test.Init(t)

	res := callEndpoint(t, initResult.Server, map[string]any{
		"title":          "",
		"body":           "",
		"isPublished":    "false",
		"publishStartAt": "invalid",
		"publishEndAt":   "",
		"categoryIds":    []string{"invalid"},
	})

	res.AssertRedirectTo(t, "/admin/article/new")
	res.AssertValidationError(t, initResult.Server, initResult.SessionManager, handlererror.ValidationErrorMessages{
		"title":          "タイトルを入力してください。",
		"body":           "本文を入力してください。",
		"publishStartAt": "日時の形式が不正です。",
		"categoryIds":    "選択したカテゴリが不正です。",
	})

	assertArticleCount(t, initResult.EntClient, 0)
}

func Test存在しないカテゴリならエラー(t *testing.T) {
	initResult := test.Init(t)

	res := callEndpoint(t, initResult.Server, map[string]any{
		"title":          "title",
		"body":           "content",
		"isPublished":    "false",
		"publishStartAt": "",
		"publishEndAt":   "",
		"categoryIds":    []string{"999"},
	})

	res.AssertRedirectTo(t, "/admin/article/new")
	res.AssertFlashError(t, initResult.Server, initResult.SessionManager, "選択したカテゴリが見つかりません。")

	assertArticleCount(t, initResult.EntClient, 0)
}

func Testカテゴリの検索に失敗した場合はエラー(t *testing.T) {
	initResult := test.Init(t, &di.ContainerOptions{
		CategoryRepository: stubCategory.CategoryRepositoryStub{
			SearchFunc: func(ctx context.Context, criteria domainCategory.SearchCategoryCriteria) ([]domainCategory.Category, error) {
				return nil, errors.New("test")
			},
		},
	})

	res := callEndpoint(t, initResult.Server, map[string]any{
		"title":          "title",
		"body":           "content",
		"isPublished":    "false",
		"publishStartAt": "",
		"publishEndAt":   "",
		"categoryIds":    []string{"1"},
	})

	res.AssertRedirectTo(t, "/admin/article/new")
	res.AssertFlashError(t, initResult.Server, initResult.SessionManager, "カテゴリの読み込みに失敗しました。")

	assertArticleCount(t, initResult.EntClient, 0)
}

func Test記事の作成に失敗した場合はエラー(t *testing.T) {
	initResult := test.Init(t, &di.ContainerOptions{
		ArticleRepository: stubArticle.ArticleRepositoryStub{
			CreateFunc: func(ctx context.Context, input domainArticle.CreateArticleInput) error {
				return errors.New("test")
			},
		},
	})

	res := callEndpoint(t, initResult.Server, map[string]any{
		"title":          "title",
		"body":           "content",
		"isPublished":    "false",
		"publishStartAt": "",
		"publishEndAt":   "",
		"categoryIds":    []string{},
	})

	res.AssertRedirectTo(t, "/admin/article/new")
	res.AssertFlashError(t, initResult.Server, initResult.SessionManager, "記事の作成に失敗しました。")

	assertArticleCount(t, initResult.EntClient, 0)
}

func callEndpoint(t *testing.T, server *httptest.Server, body map[string]any) inertiaAction.TestActionResponse {
	t.Helper()

	return inertiaAction.Send(t, server, inertiaAction.TestActionRequest{
		Method:       http.MethodPost,
		Path:         "/admin/article/new",
		Body:         body,
		UseBasicAuth: true,
		Referer:      "/admin/article/new",
	})
}

func fetchArticles(t *testing.T, entClient *ent.Client) []*ent.Article {
	t.Helper()

	articles, err := entClient.Article.Query().
		WithCategories().
		Order(entArticle.ByID()).
		All(context.Background())
	if err != nil {
		t.Fatalf("記事の取得に失敗: %v", err)
	}
	return articles
}

func assertArticleCount(t *testing.T, entClient *ent.Client, expected int) {
	t.Helper()

	articles := fetchArticles(t, entClient)
	helper.AssertEqual(t, expected, len(articles), "記事件数が不正です")
}

func getArticleCategoryIDs(article *ent.Article) []uint32 {
	categoryIDs := make([]uint32, 0, len(article.Edges.Categories))
	for _, category := range article.Edges.Categories {
		categoryIDs = append(categoryIDs, category.ID)
	}
	return categoryIDs
}

func categoryIDString(categoryID uint32) string {
	return strconv.FormatUint(uint64(categoryID), 10)
}
