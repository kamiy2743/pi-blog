package store_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"blog/internal/db/ent"
	"blog/internal/di"
	domainCategory "blog/internal/domain/category"
	"blog/internal/test"
	fixtureCategory "blog/internal/test/fixture/category"
	"blog/internal/test/helper"
	inertiaAction "blog/internal/test/helper/inertia/action"
	stubCategory "blog/internal/test/stub/category"
)

func Testカテゴリを作成できる(t *testing.T) {
	initResult := test.Init(t)

	res := callEndpoint(t, initResult.Server, "Go")
	res.AssertRedirectTo(t, "/admin/category")

	categories := fetchCategories(t, initResult.EntClient)
	assertCategoryCount(t, initResult.EntClient, 1)
	helper.AssertEqual(t, "Go", categories[0].Name, "カテゴリ名が不正です")
}

func Testカテゴリ名が空ならバリデーションエラー(t *testing.T) {
	initResult := test.Init(t)

	res := callEndpoint(t, initResult.Server, "")
	res.AssertRedirectTo(t, "/admin/category")

	res.AssertValidationError(t, initResult.Server, initResult.SessionManager, map[string]string{
		"create.name": "カテゴリ名を入力してください。",
	})

	assertCategoryCount(t, initResult.EntClient, 0)
}

func Testカテゴリ名が64文字を超えたらバリデーションエラー(t *testing.T) {
	initResult := test.Init(t)

	res := callEndpoint(t, initResult.Server, helper.StringOfLength(65))
	res.AssertRedirectTo(t, "/admin/category")

	res.AssertValidationError(t, initResult.Server, initResult.SessionManager, map[string]string{
		"create.name": "カテゴリ名は64文字以下で入力してください。",
	})

	assertCategoryCount(t, initResult.EntClient, 0)
}

func Testカテゴリ名が重複した場合はエラー(t *testing.T) {
	initResult := test.Init(t)
	fixtureCategory.CreateCategory(t, initResult.EntClient, fixtureCategory.CreateCategoryInput{Name: "Go"})

	res := callEndpoint(t, initResult.Server, "Go")
	res.AssertRedirectTo(t, "/admin/category")
	res.AssertFlashError(t, initResult.Server, initResult.SessionManager, "同じ名前のカテゴリが既に存在しています。")

	assertCategoryCount(t, initResult.EntClient, 1)
}

func Testカテゴリの検索に失敗した場合はエラー(t *testing.T) {
	initResult := test.Init(t, &di.ContainerOptions{
		CategoryRepository: stubCategory.CategoryRepositoryStub{
			SearchFunc: func(ctx context.Context, criteria domainCategory.SearchCategoryCriteria) ([]domainCategory.Category, error) {
				return nil, errors.New("test")
			},
		},
	})

	res := callEndpoint(t, initResult.Server, "Go")
	res.AssertRedirectTo(t, "/admin/category")
	res.AssertFlashError(t, initResult.Server, initResult.SessionManager, "カテゴリの作成に失敗しました。")

	assertCategoryCount(t, initResult.EntClient, 0)
}

func Testカテゴリの作成に失敗した場合はエラー(t *testing.T) {
	initResult := test.Init(t, &di.ContainerOptions{
		CategoryRepository: stubCategory.CategoryRepositoryStub{
			CreateFunc: func(ctx context.Context, input domainCategory.CreateCategoryInput) error {
				return errors.New("test")
			},
		},
	})

	res := callEndpoint(t, initResult.Server, "Go")
	res.AssertRedirectTo(t, "/admin/category")
	res.AssertFlashError(t, initResult.Server, initResult.SessionManager, "カテゴリの作成に失敗しました。")

	assertCategoryCount(t, initResult.EntClient, 0)
}

func callEndpoint(t *testing.T, server *httptest.Server, name string) inertiaAction.TestActionResponse {
	t.Helper()

	return inertiaAction.Send(t, server, inertiaAction.TestActionRequest{
		Method: http.MethodPost,
		Path:   "/admin/category",
		Body: map[string]any{
			"name": name,
		},
		UseBasicAuth: true,
		Referer:      "/admin/category",
	})
}

func fetchCategories(t *testing.T, entClient *ent.Client) []*ent.Category {
	t.Helper()

	categories, err := entClient.Category.Query().All(context.Background())
	if err != nil {
		t.Fatalf("カテゴリの取得に失敗: %v", err)
	}
	return categories
}

func assertCategoryCount(t *testing.T, entClient *ent.Client, expected int) {
	t.Helper()

	categories := fetchCategories(t, entClient)
	helper.AssertEqual(t, expected, len(categories), "カテゴリ件数が不正です")
}
