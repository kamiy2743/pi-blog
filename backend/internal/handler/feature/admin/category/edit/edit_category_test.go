package edit_test

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
	stubCategory "blog/internal/test/stub/category"
)

func Testカテゴリ編集画面を表示できる(t *testing.T) {
	initResult := test.Init(t)
	categories := setUpRecords(t, initResult.EntClient)

	res := callEndpoint(t, initResult.Server)

	res.AssertFullProps(t, "admin/EditCategory", map[string]any{
		"categories": []map[string]any{
			{
				"id":   categories[0].ID,
				"name": "Docker",
			},
			{
				"id":   categories[1].ID,
				"name": "Go",
			},
		},
	})
}

func Testカテゴリの取得に失敗した場合は500(t *testing.T) {
	initResult := test.Init(t, &di.ContainerOptions{
		CategoryRepository: stubCategory.CategoryRepositoryStub{
			AllFunc: func(ctx context.Context, orderBy domainCategory.OrderBy) ([]domainCategory.Category, error) {
				return nil, errors.New("test")
			},
		},
	})

	res := callEndpoint(t, initResult.Server)

	res.AssertError(t, 500, "カテゴリの読み込みに失敗しました。", "時間をおいてから、もう一度お試しください。")
}

func setUpRecords(t *testing.T, entClient *ent.Client) []*ent.Category {
	t.Helper()

	categoryDocker := fixtureCategory.CreateCategory(t, entClient, fixtureCategory.CreateCategoryInput{Name: "Docker"})
	categoryGo := fixtureCategory.CreateCategory(t, entClient, fixtureCategory.CreateCategoryInput{Name: "Go"})

	return []*ent.Category{
		categoryDocker,
		categoryGo,
	}
}

func callEndpoint(t *testing.T, server *httptest.Server) helper.TestInertiaResponse {
	t.Helper()

	return helper.RequestInertia(t, server, helper.TestInertiaRequest{
		Method:       http.MethodGet,
		Path:         "/admin/category",
		UseBasicAuth: true,
	})
}
