package search_test

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
)

type records struct {
	CategoryCloudflare  *ent.Category
	CategoryDocker      *ent.Category
	CategoryGo          *ent.Category
	CategoryRaspberryPi *ent.Category
	Articles            []*ent.Article
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

	res.AssertPartialProps(t, "article/ShowArticleList", "initial", map[string]any{
		"categories": []map[string]any{
			{
				"id":   records.CategoryCloudflare.ID,
				"name": "Cloudflare",
			},
			{
				"id":   records.CategoryDocker.ID,
				"name": "Docker",
			},
			{
				"id":   records.CategoryGo.ID,
				"name": "Go",
			},
			{
				"id":   records.CategoryRaspberryPi.ID,
				"name": "Raspberry Pi",
			},
		},
	})
}

func Testパラメータなしで記事一覧を表示できる(t *testing.T) {
	initResult := test.Init(t)
	records := setUpRecords(t, initResult.EntClient, 4)

	res := callEndpoint(t, initResult.Server, queryParams{})

	res.AssertPartialProps(t, "article/ShowArticleList", "partialSearch", map[string]any{
		"articles": []map[string]any{
			{
				"id":    records.Articles[3].ID,
				"title": "Docker Raspberry Pi",
				"date":  "2026-01-04T00:00:00Z",
				"categoryNames": []string{
					"Docker",
					"Raspberry Pi",
				},
			},
			{
				"id":    records.Articles[2].ID,
				"title": "Go Cloudflare",
				"date":  "2026-01-03T00:00:00Z",
				"categoryNames": []string{
					"Cloudflare",
					"Go",
				},
			},
			{
				"id":            records.Articles[1].ID,
				"title":         "no category",
				"date":          "2026-01-02T00:00:00Z",
				"categoryNames": []string{},
			},
			{
				"id":    records.Articles[0].ID,
				"title": "all categories",
				"date":  "2026-01-01T00:00:00Z",
				"categoryNames": []string{
					"Cloudflare",
					"Docker",
					"Go",
					"Raspberry Pi",
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
		Title: "cloud",
		CategoryIDs: []string{
			strconv.FormatUint(uint64(records.CategoryCloudflare.ID), 10),
		},
		Page: "1",
	})

	res.AssertPartialProps(t, "article/ShowArticleList", "partialSearch", map[string]any{
		"articles": []map[string]any{
			{
				"id":    records.Articles[2].ID,
				"title": "Go Cloudflare",
				"date":  "2026-01-03T00:00:00Z",
				"categoryNames": []string{
					"Cloudflare",
					"Go",
				},
			},
		},
		"title": "cloud",
		"categoryIds": []uint32{
			records.CategoryCloudflare.ID,
		},
		"page":       1,
		"totalCount": 1,
		"totalPages": 1,
	})
}

func Testページネーションできる(t *testing.T) {
	initResult := test.Init(t)
	setUpRecords(t, initResult.EntClient, 33)

	res := callEndpoint(t, initResult.Server, queryParams{
		Page: "4",
	})

	res.AssertPropsCount(t, "article/ShowArticleList", "partialSearch.articles", 3)
	res.AssertPropsValue(t, "article/ShowArticleList", "partialSearch.page", 4)
	res.AssertPropsValue(t, "article/ShowArticleList", "partialSearch.totalCount", 33)
	res.AssertPropsValue(t, "article/ShowArticleList", "partialSearch.totalPages", 4)
}

func TestカテゴリIDが不正なら除外する(t *testing.T) {
	initResult := test.Init(t)
	records := setUpRecords(t, initResult.EntClient, 4)

	res := callEndpoint(t, initResult.Server, queryParams{
		CategoryIDs: []string{
			strconv.FormatUint(uint64(records.CategoryCloudflare.ID), 10),
			"invalid",
			strconv.FormatUint(uint64(99), 10),
		},
	})

	res.AssertPropsValue(t, "article/ShowArticleList", "partialSearch.categoryIds", []uint32{
		records.CategoryCloudflare.ID,
	})
}

func Testページ番号が不正なら1ページ目を表示する(t *testing.T) {
	initResult := test.Init(t)
	setUpRecords(t, initResult.EntClient, 4)

	res := callEndpoint(t, initResult.Server, queryParams{
		Page: "invalid",
	})

	res.AssertPropsValue(t, "article/ShowArticleList", "partialSearch.page", 1)
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

func setUpRecords(
	t *testing.T,
	entClient *ent.Client,
	requireTotalArticleCount int,
) records {
	t.Helper()

	categoryCloudflare := fixtureCategory.CreateCategory(t, entClient, fixtureCategory.CreateCategoryInput{Name: "Cloudflare"})
	categoryDocker := fixtureCategory.CreateCategory(t, entClient, fixtureCategory.CreateCategoryInput{Name: "Docker"})
	categoryGo := fixtureCategory.CreateCategory(t, entClient, fixtureCategory.CreateCategoryInput{Name: "Go"})
	categoryRaspberryPi := fixtureCategory.CreateCategory(t, entClient, fixtureCategory.CreateCategoryInput{Name: "Raspberry Pi"})

	articles := make([]*ent.Article, 0, requireTotalArticleCount)
	articles = append(articles, fixtureArticle.CreateArticle(
		t,
		entClient,
		fixtureArticle.CreateArticleInput{
			Title:       "all categories",
			Body:        "content",
			IsPublished: true,
			UpdatedAt:   helper.TimePtr(t, "2026-01-01 00:00"),
			Categories: []*ent.Category{
				categoryCloudflare,
				categoryDocker,
				categoryGo,
				categoryRaspberryPi,
			},
		},
	))
	articles = append(articles, fixtureArticle.CreateArticle(
		t,
		entClient,
		fixtureArticle.CreateArticleInput{
			Title:       "no category",
			Body:        "content",
			IsPublished: true,
			UpdatedAt:   helper.TimePtr(t, "2026-01-02 00:00"),
			Categories:  []*ent.Category{},
		},
	))
	articles = append(articles, fixtureArticle.CreateArticle(
		t,
		entClient,
		fixtureArticle.CreateArticleInput{
			Title:       "Go Cloudflare",
			Body:        "content",
			IsPublished: true,
			UpdatedAt:   helper.TimePtr(t, "2026-01-03 00:00"),
			Categories: []*ent.Category{
				categoryCloudflare,
				categoryGo,
			},
		},
	))
	articles = append(articles, fixtureArticle.CreateArticle(
		t,
		entClient,
		fixtureArticle.CreateArticleInput{
			Title:       "Docker Raspberry Pi",
			Body:        "content",
			IsPublished: true,
			UpdatedAt:   helper.TimePtr(t, "2026-01-04 00:00"),
			Categories: []*ent.Category{
				categoryDocker,
				categoryRaspberryPi,
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
				IsPublished: true,
			},
		))
	}

	return records{
		CategoryCloudflare:  categoryCloudflare,
		CategoryDocker:      categoryDocker,
		CategoryGo:          categoryGo,
		CategoryRaspberryPi: categoryRaspberryPi,
		Articles:            articles,
	}
}

func callEndpoint(
	t *testing.T,
	server *httptest.Server,
	params queryParams,
) helper.TestInertiaResponse {
	t.Helper()

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

	return helper.RequestInertia(t, server, helper.TestInertiaRequest{
		Method:      http.MethodGet,
		Path:        "/article",
		QueryParams: query,
	})
}
