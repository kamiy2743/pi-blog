package show_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"blog/internal/db/ent"
	"blog/internal/di"
	domainArticle "blog/internal/domain/article"
	domainCategory "blog/internal/domain/category"
	"blog/internal/test"
	fixtureArticle "blog/internal/test/fixture/article"
	fixtureCategory "blog/internal/test/fixture/category"
	"blog/internal/test/helper"
	stubArticle "blog/internal/test/stub/article"
	stubCategory "blog/internal/test/stub/category"

	"github.com/romsar/gonertia/v3"
)

type records struct {
	CategoryDocker *ent.Category
	CategoryGo     *ent.Category
	Articles       []*ent.Article
}

type queryParams struct {
	Title       string
	CategoryIDs []string
	Page        string
}

func Testカテゴリ一覧を表示できる(t *testing.T) {
	initResult := test.Init(t)
	records := setUpRecords(t, initResult.EntClient, 4)

	res := callEndpoint(t, initResult.Server, queryParams{})

	res.AssertPartialProps(t, "admin/ShowAdmin", "initial", gonertia.Props{
		"categories": []gonertia.Props{
			{
				"id":   records.CategoryDocker.ID,
				"name": "Docker",
			},
			{
				"id":   records.CategoryGo.ID,
				"name": "Go",
			},
		},
	})
}

func Testパラメータなしで記事一覧を表示できる(t *testing.T) {
	initResult := test.Init(t)
	records := setUpRecords(t, initResult.EntClient, 4)

	res := callEndpoint(t, initResult.Server, queryParams{})

	res.AssertPartialProps(t, "admin/ShowAdmin", "partialSearch", gonertia.Props{
		"articles": []gonertia.Props{
			{
				"id":          records.Articles[3].ID,
				"title":       "Docker 非公開",
				"date":        "2026-01-04T00:00:00Z",
				"isPublished": false,
				"categoryNames": []string{
					"Docker",
				},
			},
			{
				"id":          records.Articles[2].ID,
				"title":       "Go 公開",
				"date":        "2026-01-03T00:00:00Z",
				"isPublished": true,
				"categoryNames": []string{
					"Go",
				},
			},
			{
				"id":            records.Articles[1].ID,
				"title":         "no category",
				"date":          "2026-01-02T00:00:00Z",
				"isPublished":   false,
				"categoryNames": []string{},
			},
			{
				"id":          records.Articles[0].ID,
				"title":       "Go Docker 公開",
				"date":        "2026-01-01T00:00:00Z",
				"isPublished": true,
				"categoryNames": []string{
					"Docker",
					"Go",
				},
			},
		},
		"title":       "",
		"categoryIds": []uint32{},
		"page":        1,
		"totalCount":  4,
		"totalPages":  1,
	})
}

func Test全パラメータありで記事一覧を表示できる(t *testing.T) {
	initResult := test.Init(t)
	records := setUpRecords(t, initResult.EntClient, 4)

	res := callEndpoint(t, initResult.Server, queryParams{
		Title: "docker",
		CategoryIDs: []string{
			strconv.FormatUint(uint64(records.CategoryDocker.ID), 10),
		},
		Page: "1",
	})

	res.AssertPartialProps(t, "admin/ShowAdmin", "partialSearch", gonertia.Props{
		"articles": []gonertia.Props{
			{
				"id":          records.Articles[3].ID,
				"title":       "Docker 非公開",
				"date":        "2026-01-04T00:00:00Z",
				"isPublished": false,
				"categoryNames": []string{
					"Docker",
				},
			},
			{
				"id":          records.Articles[0].ID,
				"title":       "Go Docker 公開",
				"date":        "2026-01-01T00:00:00Z",
				"isPublished": true,
				"categoryNames": []string{
					"Docker",
					"Go",
				},
			},
		},
		"title": "docker",
		"categoryIds": []uint32{
			records.CategoryDocker.ID,
		},
		"page":       1,
		"totalCount": 2,
		"totalPages": 1,
	})
}

func Testページネーションできる(t *testing.T) {
	initResult := test.Init(t)
	setUpRecords(t, initResult.EntClient, 33)

	res := callEndpoint(t, initResult.Server, queryParams{
		Page: "2",
	})

	res.AssertPropsCount(t, "admin/ShowAdmin", "partialSearch.articles", 13)
	res.AssertPropsValue(t, "admin/ShowAdmin", "partialSearch.page", 2)
	res.AssertPropsValue(t, "admin/ShowAdmin", "partialSearch.totalCount", 33)
	res.AssertPropsValue(t, "admin/ShowAdmin", "partialSearch.totalPages", 2)
}

func TestカテゴリIDが不正なら除外する(t *testing.T) {
	initResult := test.Init(t)
	records := setUpRecords(t, initResult.EntClient, 4)

	res := callEndpoint(t, initResult.Server, queryParams{
		CategoryIDs: []string{
			strconv.FormatUint(uint64(records.CategoryDocker.ID), 10),
			"invalid",
			strconv.FormatUint(uint64(99), 10),
		},
	})

	res.AssertPropsValue(t, "admin/ShowAdmin", "partialSearch.categoryIds", []uint32{
		records.CategoryDocker.ID,
	})
}

func Testページ番号が不正なら1ページ目を表示する(t *testing.T) {
	initResult := test.Init(t)
	setUpRecords(t, initResult.EntClient, 4)

	res := callEndpoint(t, initResult.Server, queryParams{
		Page: "invalid",
	})

	res.AssertPropsValue(t, "admin/ShowAdmin", "partialSearch.page", 1)
}

func Testタイトルが長すぎる場合はバリデーションエラー(t *testing.T) {
	initResult := test.Init(t)
	records := setUpRecords(t, initResult.EntClient, 4)

	title := helper.StringOfLength(256)
	res := callEndpoint(t, initResult.Server, queryParams{
		Title: title,
		CategoryIDs: []string{
			strconv.FormatUint(uint64(records.CategoryDocker.ID), 10),
		},
	})

	res.AssertPropsValue(t, "admin/ShowAdmin", "validationErrors.title", "タイトルは255文字以下で入力してください。")
	res.AssertPartialProps(t, "admin/ShowAdmin", "partialSearch", gonertia.Props{
		"title": title,
		"categoryIds": []uint32{
			records.CategoryDocker.ID,
		},
		"page":       1,
		"totalCount": 0,
		"totalPages": 1,
		"articles":   []gonertia.Props{},
	})
}

func TestPartialReloadでタイトルが長すぎる場合はバリデーションエラー(t *testing.T) {
	initResult := test.Init(t)
	records := setUpRecords(t, initResult.EntClient, 4)

	title := helper.StringOfLength(256)
	res := callPartialEndpoint(t, initResult.Server, queryParams{
		Title: title,
		CategoryIDs: []string{
			strconv.FormatUint(uint64(records.CategoryDocker.ID), 10),
		},
	})

	res.AssertPropsValue(t, "admin/ShowAdmin", "validationErrors.title", "タイトルは255文字以下で入力してください。")
	res.AssertPartialProps(t, "admin/ShowAdmin", "partialSearch", gonertia.Props{
		"title": title,
		"categoryIds": []uint32{
			records.CategoryDocker.ID,
		},
		"page":       1,
		"totalCount": 0,
		"totalPages": 1,
		"articles":   []gonertia.Props{},
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
	setUpRecords(t, initResult.EntClient, 4)

	res := callEndpoint(t, initResult.Server, queryParams{})

	res.AssertError(t, 500, "カテゴリの読み込みに失敗しました。", "時間をおいてから、もう一度お試しください。")
}

func Test記事の取得に失敗した場合は500(t *testing.T) {
	initResult := test.Init(t, &di.ContainerOptions{
		ArticleRepository: stubArticle.ArticleRepositoryStub{
			PaginateFunc: func(ctx context.Context, criteria domainArticle.PaginateArticleCriteria) (domainArticle.PaginatedArticles, error) {
				return domainArticle.PaginatedArticles{}, errors.New("test")
			},
		},
	})
	setUpRecords(t, initResult.EntClient, 4)

	res := callEndpoint(t, initResult.Server, queryParams{})

	res.AssertError(t, 500, "記事の読み込みに失敗しました。", "時間をおいてから、もう一度お試しください。")
}

func TestPartialReloadで記事の取得に失敗した場合はフラッシュメッセージを表示する(t *testing.T) {
	initResult := test.Init(t, &di.ContainerOptions{
		ArticleRepository: stubArticle.ArticleRepositoryStub{
			PaginateFunc: func(ctx context.Context, criteria domainArticle.PaginateArticleCriteria) (domainArticle.PaginatedArticles, error) {
				return domainArticle.PaginatedArticles{}, errors.New("test")
			},
		},
	})
	setUpRecords(t, initResult.EntClient, 4)

	res := callPartialEndpoint(t, initResult.Server, queryParams{
		Title: "docker",
		Page:  "2",
	})

	res.AssertPropsValue(t, "admin/ShowAdmin", "flash.error", "記事の読み込みに失敗しました。")
	res.AssertPartialProps(t, "admin/ShowAdmin", "partialSearch", gonertia.Props{
		"title":       "docker",
		"categoryIds": []uint32{},
		"page":        1,
		"totalCount":  0,
		"totalPages":  1,
		"articles":    []gonertia.Props{},
	})
}

func setUpRecords(
	t *testing.T,
	entClient *ent.Client,
	requireTotalArticleCount int,
) records {
	t.Helper()

	categoryDocker := fixtureCategory.CreateCategory(t, entClient, fixtureCategory.CreateCategoryInput{Name: "Docker"})
	categoryGo := fixtureCategory.CreateCategory(t, entClient, fixtureCategory.CreateCategoryInput{Name: "Go"})

	articles := make([]*ent.Article, 0, requireTotalArticleCount)
	articles = append(articles, fixtureArticle.CreateArticle(
		t,
		entClient,
		fixtureArticle.CreateArticleInput{
			Title:       "Go Docker 公開",
			Body:        "content",
			IsPublished: true,
			UpdatedAt:   helper.TimePtr(t, "2026-01-01 00:00"),
			Categories: []*ent.Category{
				categoryDocker,
				categoryGo,
			},
		},
	))
	articles = append(articles, fixtureArticle.CreateArticle(
		t,
		entClient,
		fixtureArticle.CreateArticleInput{
			Title:       "no category",
			Body:        "content",
			IsPublished: false,
			UpdatedAt:   helper.TimePtr(t, "2026-01-02 00:00"),
			Categories:  []*ent.Category{},
		},
	))
	articles = append(articles, fixtureArticle.CreateArticle(
		t,
		entClient,
		fixtureArticle.CreateArticleInput{
			Title:       "Go 公開",
			Body:        "content",
			IsPublished: true,
			UpdatedAt:   helper.TimePtr(t, "2026-01-03 00:00"),
			Categories: []*ent.Category{
				categoryGo,
			},
		},
	))
	articles = append(articles, fixtureArticle.CreateArticle(
		t,
		entClient,
		fixtureArticle.CreateArticleInput{
			Title:       "Docker 非公開",
			Body:        "content",
			IsPublished: false,
			UpdatedAt:   helper.TimePtr(t, "2026-01-04 00:00"),
			Categories: []*ent.Category{
				categoryDocker,
			},
		},
	))
	presetArticleCount := 4

	if requireTotalArticleCount < presetArticleCount {
		t.Fatalf("requireTotalArticleCountは%d以上を指定してください。現在の値: %d", presetArticleCount, requireTotalArticleCount)
	}

	for i := 1; i <= requireTotalArticleCount-presetArticleCount; i++ {
		articles = append(articles, fixtureArticle.CreateArticle(
			t,
			entClient,
			fixtureArticle.CreateArticleInput{
				IsPublished: i%2 == 0,
			},
		))
	}

	return records{
		CategoryDocker: categoryDocker,
		CategoryGo:     categoryGo,
		Articles:       articles,
	}
}

func callEndpoint(
	t *testing.T,
	server *httptest.Server,
	params queryParams,
) helper.TestInertiaResponse {
	t.Helper()

	return helper.RequestInertia(t, server, helper.TestInertiaRequest{
		Method:       http.MethodGet,
		Path:         "/admin",
		QueryParams:  buildQuery(params),
		UseBasicAuth: true,
	})
}

func callPartialEndpoint(
	t *testing.T,
	server *httptest.Server,
	params queryParams,
) helper.TestInertiaResponse {
	t.Helper()

	return helper.RequestInertia(t, server, helper.TestInertiaRequest{
		Method:           http.MethodGet,
		Path:             "/admin",
		QueryParams:      buildQuery(params),
		UseBasicAuth:     true,
		PartialComponent: "admin/ShowAdmin",
		PartialData:      []string{"partialSearch", "validationErrors", "flash"},
	})
}

func buildQuery(params queryParams) map[string][]string {
	query := map[string][]string{}

	if params.Title != "" {
		query["title"] = []string{params.Title}
	}
	if len(params.CategoryIDs) > 0 {
		query["categoryId"] = params.CategoryIDs
	}
	if params.Page != "" {
		query["page"] = []string{params.Page}
	}

	return query
}
