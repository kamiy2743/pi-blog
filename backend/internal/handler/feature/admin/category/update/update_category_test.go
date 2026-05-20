package update_test

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
	"blog/internal/handler/handlererror"
	"blog/internal/handler/session"
	"blog/internal/test"
	fixtureCategory "blog/internal/test/fixture/category"
	"blog/internal/test/helper"
	inertiaAction "blog/internal/test/helper/inertia/action"
	stubCategory "blog/internal/test/stub/category"
)

func Testカテゴリを更新できる(t *testing.T) {
	initResult := test.Init(t)
	category := setUpRecord(t, initResult.EntClient)

	res := callEndpoint(t, initResult.Server, category.ID, "Docker")
	res.AssertRedirectTo(t, "/admin/category")

	categories := fetchCategories(t, initResult.EntClient)
	assertCategoryCount(t, initResult.EntClient, 1)
	helper.AssertEqual(t, "Docker", categories[0].Name, "カテゴリ名が不正です")
}

func Testカテゴリ名が空ならバリデーションエラー(t *testing.T) {
	initResult := test.Init(t)
	category := setUpRecord(t, initResult.EntClient)

	res := callEndpoint(t, initResult.Server, category.ID, "")
	res.AssertRedirectTo(t, "/admin/category")

	res.AssertOldInput(t, initResult.Server, initResult.SessionManager, session.OldInput{
		"formKey": "category.update." + strconv.FormatUint(uint64(category.ID), 10),
		"name":    "",
	})
	res.AssertValidationError(t, initResult.Server, initResult.SessionManager, handlererror.ValidationErrorMessages{
		"name": "カテゴリ名を入力してください。",
	})

	categories := fetchCategories(t, initResult.EntClient)
	helper.AssertEqual(t, "Go", categories[0].Name, "カテゴリ名が更新されてはいけません")
}

func Testカテゴリ名が64文字を超えたらバリデーションエラー(t *testing.T) {
	initResult := test.Init(t)
	category := setUpRecord(t, initResult.EntClient)

	name := helper.StringOfLength(65)
	res := callEndpoint(t, initResult.Server, category.ID, name)
	res.AssertRedirectTo(t, "/admin/category")

	res.AssertOldInput(t, initResult.Server, initResult.SessionManager, session.OldInput{
		"formKey": "category.update." + strconv.FormatUint(uint64(category.ID), 10),
		"name":    name,
	})
	res.AssertValidationError(t, initResult.Server, initResult.SessionManager, handlererror.ValidationErrorMessages{
		"name": "カテゴリ名は64文字以下で入力してください。",
	})

	categories := fetchCategories(t, initResult.EntClient)
	helper.AssertEqual(t, "Go", categories[0].Name, "カテゴリ名が更新されてはいけません")
}

func Testカテゴリ名が重複した場合はエラー(t *testing.T) {
	initResult := test.Init(t)
	setUpRecord(t, initResult.EntClient)
	category := fixtureCategory.CreateCategory(t, initResult.EntClient, fixtureCategory.CreateCategoryInput{Name: "Docker"})

	res := callEndpoint(t, initResult.Server, category.ID, "Go")
	res.AssertRedirectTo(t, "/admin/category")

	res.AssertOldInput(t, initResult.Server, initResult.SessionManager, session.OldInput{
		"formKey": "category.update." + strconv.FormatUint(uint64(category.ID), 10),
		"name":    "Go",
	})
	res.AssertFlashError(t, initResult.Server, initResult.SessionManager, "同じ名前のカテゴリが既に存在しています。")

	categories := fetchCategories(t, initResult.EntClient)
	helper.AssertEqual(t, "Docker", categories[1].Name, "カテゴリ名が更新されてはいけません")
}

func Testカテゴリが存在しない場合はエラー(t *testing.T) {
	initResult := test.Init(t)

	res := callEndpoint(t, initResult.Server, 1, "Go")
	res.AssertRedirectTo(t, "/admin/category")
	res.AssertFlashError(t, initResult.Server, initResult.SessionManager, "カテゴリが見つかりません。")

	assertCategoryCount(t, initResult.EntClient, 0)
}

func Testカテゴリの検索に失敗した場合はエラー(t *testing.T) {
	initResult := test.Init(t, &di.ContainerOptions{
		CategoryRepository: stubCategory.CategoryRepositoryStub{
			SearchFunc: func(ctx context.Context, criteria domainCategory.SearchCategoryCriteria) ([]domainCategory.Category, error) {
				return nil, errors.New("test")
			},
		},
	})

	res := callEndpoint(t, initResult.Server, 1, "Go")
	res.AssertRedirectTo(t, "/admin/category")
	res.AssertFlashError(t, initResult.Server, initResult.SessionManager, "カテゴリの更新に失敗しました。")
}

func Testカテゴリの更新に失敗した場合はエラー(t *testing.T) {
	initResult := test.Init(t, &di.ContainerOptions{
		CategoryRepository: stubCategory.CategoryRepositoryStub{
			SearchFunc: func(ctx context.Context, criteria domainCategory.SearchCategoryCriteria) ([]domainCategory.Category, error) {
				if len(criteria.IDs) > 0 {
					return []domainCategory.Category{
						{ID: 1, Name: "Go"},
					}, nil
				}
				return []domainCategory.Category{}, nil
			},
			UpdateFunc: func(ctx context.Context, category domainCategory.Category) error {
				return errors.New("test")
			},
		},
	})

	res := callEndpoint(t, initResult.Server, 1, "Docker")
	res.AssertRedirectTo(t, "/admin/category")
	res.AssertFlashError(t, initResult.Server, initResult.SessionManager, "カテゴリの更新に失敗しました。")
}

func callEndpoint(t *testing.T, server *httptest.Server, categoryID uint32, name string) inertiaAction.TestActionResponse {
	t.Helper()

	return inertiaAction.Send(t, server, inertiaAction.TestActionRequest{
		Method: http.MethodPost,
		Path:   "/admin/category/" + strconv.FormatUint(uint64(categoryID), 10),
		Body: map[string]any{
			"formKey": "category.update." + strconv.FormatUint(uint64(categoryID), 10),
			"name":    name,
		},
		UseBasicAuth: true,
		Referer:      "/admin/category",
	})
}

func setUpRecord(t *testing.T, entClient *ent.Client) *ent.Category {
	t.Helper()

	return fixtureCategory.CreateCategory(t, entClient, fixtureCategory.CreateCategoryInput{Name: "Go"})
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
