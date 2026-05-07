package destroy_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"blog/internal/db/ent"
	"blog/internal/di"
	domainCategory "blog/internal/domain/category"
	"blog/internal/test"
	fixtureArticle "blog/internal/test/fixture/article"
	fixtureCategory "blog/internal/test/fixture/category"
	"blog/internal/test/helper"
	inertiaAction "blog/internal/test/helper/inertia/action"
	stubCategory "blog/internal/test/stub/category"
)

func Testカテゴリを削除できる(t *testing.T) {
	initResult := test.Init(t)
	category := setUpRecord(t, initResult.EntClient)

	res := callEndpoint(t, initResult.Server, category.ID)
	res.AssertRedirectTo(t, "/admin/category")

	assertCategoryCount(t, initResult.EntClient, 0)
}

func Test記事に紐づくカテゴリを削除できる(t *testing.T) {
	initResult := test.Init(t)
	category := setUpRecord(t, initResult.EntClient)
	fixtureArticle.CreateArticle(t, initResult.EntClient, fixtureArticle.CreateArticleInput{
		Title:      "Go",
		Body:       "本文",
		Categories: []*ent.Category{category},
	})

	res := callEndpoint(t, initResult.Server, category.ID)
	res.AssertRedirectTo(t, "/admin/category")

	assertCategoryCount(t, initResult.EntClient, 0)
	assertArticleCount(t, initResult.EntClient, 1)
}

func Testカテゴリが存在しない場合はエラー(t *testing.T) {
	initResult := test.Init(t)

	res := callEndpoint(t, initResult.Server, 1)
	res.AssertRedirectTo(t, "/admin/category")
	res.AssertFlashError(t, initResult.Server, initResult.SessionManager, "カテゴリが見つかりません。")

	assertCategoryCount(t, initResult.EntClient, 0)
}

func TestカテゴリIDが不正な場合はエラー(t *testing.T) {
	initResult := test.Init(t)

	res := inertiaAction.Send(t, initResult.Server, inertiaAction.TestActionRequest{
		Method:       http.MethodPost,
		Path:         "/admin/category/invalid/delete",
		Body:         map[string]any{},
		UseBasicAuth: true,
		Referer:      "/admin/category",
	})
	res.AssertRedirectTo(t, "/admin/category")
	res.AssertFlashError(t, initResult.Server, initResult.SessionManager, "カテゴリが見つかりません。")
}

func Testカテゴリの検索に失敗した場合はエラー(t *testing.T) {
	initResult := test.Init(t, &di.ContainerOptions{
		CategoryRepository: stubCategory.CategoryRepositoryStub{
			SearchFunc: func(ctx context.Context, criteria domainCategory.SearchCategoryCriteria) ([]domainCategory.Category, error) {
				return nil, errors.New("test")
			},
		},
	})

	res := callEndpoint(t, initResult.Server, 1)
	res.AssertRedirectTo(t, "/admin/category")
	res.AssertFlashError(t, initResult.Server, initResult.SessionManager, "カテゴリの削除に失敗しました。")
}

func Testカテゴリの削除に失敗した場合はエラー(t *testing.T) {
	initResult := test.Init(t, &di.ContainerOptions{
		CategoryRepository: stubCategory.CategoryRepositoryStub{
			SearchFunc: func(ctx context.Context, criteria domainCategory.SearchCategoryCriteria) ([]domainCategory.Category, error) {
				return []domainCategory.Category{
					{ID: 1, Name: "Go"},
				}, nil
			},
			DeleteFunc: func(ctx context.Context, category domainCategory.Category) error {
				return errors.New("test")
			},
		},
	})

	res := callEndpoint(t, initResult.Server, 1)
	res.AssertRedirectTo(t, "/admin/category")
	res.AssertFlashError(t, initResult.Server, initResult.SessionManager, "カテゴリの削除に失敗しました。")
}

func callEndpoint(t *testing.T, server *httptest.Server, categoryID uint32) inertiaAction.TestActionResponse {
	t.Helper()

	return inertiaAction.Send(t, server, inertiaAction.TestActionRequest{
		Method:       http.MethodPost,
		Path:         "/admin/category/" + strconv.FormatUint(uint64(categoryID), 10) + "/delete",
		Body:         map[string]any{},
		UseBasicAuth: true,
		Referer:      "/admin/category",
	})
}

func setUpRecord(t *testing.T, entClient *ent.Client) *ent.Category {
	t.Helper()

	return fixtureCategory.CreateCategory(t, entClient, fixtureCategory.CreateCategoryInput{Name: "Go"})
}

func assertCategoryCount(t *testing.T, entClient *ent.Client, expected int) {
	t.Helper()

	count, err := entClient.Category.Query().Count(context.Background())
	if err != nil {
		t.Fatalf("カテゴリ件数の取得に失敗: %v", err)
	}
	helper.AssertEqual(t, expected, count, "カテゴリ件数が不正です")
}

func assertArticleCount(t *testing.T, entClient *ent.Client, expected int) {
	t.Helper()

	count, err := entClient.Article.Query().Count(context.Background())
	if err != nil {
		t.Fatalf("記事件数の取得に失敗: %v", err)
	}
	helper.AssertEqual(t, expected, count, "記事件数が不正です")
}
