package update_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sort"
	"strconv"
	"testing"

	"blog/internal/db/ent"
	entArticle "blog/internal/db/ent/article"
	"blog/internal/di"
	domainArticle "blog/internal/domain/article"
	domainCategory "blog/internal/domain/category"
	"blog/internal/handler/handlererror"
	"blog/internal/handler/session"
	"blog/internal/test"
	fixtureArticle "blog/internal/test/fixture/article"
	fixtureCategory "blog/internal/test/fixture/category"
	"blog/internal/test/helper"
	inertiaAction "blog/internal/test/helper/inertia/action"
	stubArticle "blog/internal/test/stub/article"
	stubCategory "blog/internal/test/stub/category"
)

func Test記事を更新できる(t *testing.T) {
	initResult := test.Init(t)
	records := setUpRecords(t, initResult.EntClient)

	res := callEndpoint(t, initResult.Server, records.Article.ID, map[string]any{
		"title":          "new title",
		"body":           "## new",
		"isPublished":    "true",
		"publishStartAt": "2026-01-02T03:04:00Z",
		"publishEndAt":   "2026-01-03T04:05:00Z",
		"categoryIds":    []string{categoryIDString(records.CategoryGo.ID)},
	})

	res.AssertRedirectTo(t, "/admin")

	articles := fetchArticles(t, initResult.EntClient)
	helper.AssertEqual(t, 1, len(articles), "記事件数が不正です")
	helper.AssertEqual(t, "new title", articles[0].Title, "タイトルが不正です")
	helper.AssertEqual(t, "## new", articles[0].BodyMarkdown, "Markdown 本文が不正です")
	helper.AssertEqual(t, "<h2>new</h2>\n", articles[0].BodyHTML, "HTML 本文が不正です")
	helper.AssertEqual(t, true, articles[0].IsPublished, "公開状態が不正です")
	helper.AssertEqual(t, helper.TimePtr(t, "2026-01-02 03:04"), articles[0].PublishStartAt, "公開開始時刻が不正です")
	helper.AssertEqual(t, helper.TimePtr(t, "2026-01-03 04:05"), articles[0].PublishEndAt, "公開終了時刻が不正です")
	helper.AssertEqual(t, []uint32{records.CategoryGo.ID}, getArticleCategoryIDs(articles[0]), "カテゴリIDが不正です")
}

func Testタイトルが空ならバリデーションエラー(t *testing.T) {
	initResult := test.Init(t)
	records := setUpRecords(t, initResult.EntClient)

	res := callEndpoint(t, initResult.Server, records.Article.ID, map[string]any{
		"title":          "",
		"body":           "content",
		"isPublished":    "false",
		"publishStartAt": "",
		"publishEndAt":   "",
		"categoryIds":    []any{},
	})

	res.AssertRedirectTo(t, "/admin/article/"+articleIDString(records.Article.ID))
	res.AssertOldInput(t, initResult.Server, initResult.SessionManager, session.OldInput{
		"title":          "",
		"body":           "content",
		"isPublished":    "false",
		"publishStartAt": "",
		"publishEndAt":   "",
		"categoryIds":    []any{},
	})
	res.AssertValidationError(t, initResult.Server, initResult.SessionManager, handlererror.ValidationErrorMessages{
		"title": "タイトルを入力してください。",
	})
}

func Test存在しない記事ならエラー(t *testing.T) {
	initResult := test.Init(t)

	res := callEndpoint(t, initResult.Server, 999, map[string]any{
		"title":          "title",
		"body":           "content",
		"isPublished":    "false",
		"publishStartAt": "",
		"publishEndAt":   "",
		"categoryIds":    []string{},
	})

	res.AssertRedirectTo(t, "/admin/article/999")
	res.AssertFlashError(t, initResult.Server, initResult.SessionManager, "記事が見つかりません。")
}

func Test存在しないカテゴリならエラー(t *testing.T) {
	initResult := test.Init(t)
	records := setUpRecords(t, initResult.EntClient)

	res := callEndpoint(t, initResult.Server, records.Article.ID, map[string]any{
		"title":          "title",
		"body":           "content",
		"isPublished":    "false",
		"publishStartAt": "",
		"publishEndAt":   "",
		"categoryIds":    []string{"999"},
	})

	res.AssertRedirectTo(t, "/admin/article/"+articleIDString(records.Article.ID))
	res.AssertFlashError(t, initResult.Server, initResult.SessionManager, "選択したカテゴリが見つかりません。")
}

func Test記事の更新に失敗した場合はエラー(t *testing.T) {
	initResult := test.Init(t, &di.ContainerOptions{
		ArticleRepository: stubArticle.ArticleRepositoryStub{
			SearchFunc: func(ctx context.Context, criteria domainArticle.SearchArticleCriteria) ([]domainArticle.Article, error) {
				return []domainArticle.Article{
					{
						ID:           1,
						Title:        "title",
						BodyMarkdown: "body",
						BodyHTML:     "<p>body</p>\n",
						UpdatedAt:    helper.Time(t, "2026-01-02 03:04"),
					},
				}, nil
			},
			UpdateFunc: func(ctx context.Context, article domainArticle.Article) error {
				return errors.New("test")
			},
		},
	})

	res := callEndpoint(t, initResult.Server, 1, map[string]any{
		"title":          "title",
		"body":           "content",
		"isPublished":    "false",
		"publishStartAt": "",
		"publishEndAt":   "",
		"categoryIds":    []string{},
	})

	res.AssertRedirectTo(t, "/admin/article/1")
	res.AssertFlashError(t, initResult.Server, initResult.SessionManager, "記事の更新に失敗しました。")
}

func Testカテゴリの検索に失敗した場合はエラー(t *testing.T) {
	initResult := test.Init(t, &di.ContainerOptions{
		ArticleRepository: stubArticle.ArticleRepositoryStub{
			SearchFunc: func(ctx context.Context, criteria domainArticle.SearchArticleCriteria) ([]domainArticle.Article, error) {
				return []domainArticle.Article{
					{
						ID:           1,
						Title:        "title",
						BodyMarkdown: "body",
						BodyHTML:     "<p>body</p>\n",
						UpdatedAt:    helper.Time(t, "2026-01-02 03:04"),
					},
				}, nil
			},
		},
		CategoryRepository: stubCategory.CategoryRepositoryStub{
			SearchFunc: func(ctx context.Context, criteria domainCategory.SearchCategoryCriteria) ([]domainCategory.Category, error) {
				return nil, errors.New("test")
			},
		},
	})

	res := callEndpoint(t, initResult.Server, 1, map[string]any{
		"title":          "title",
		"body":           "content",
		"isPublished":    "false",
		"publishStartAt": "",
		"publishEndAt":   "",
		"categoryIds":    []string{"1"},
	})

	res.AssertRedirectTo(t, "/admin/article/1")
	res.AssertFlashError(t, initResult.Server, initResult.SessionManager, "カテゴリの読み込みに失敗しました。")
}

type records struct {
	CategoryDocker *ent.Category
	CategoryGo     *ent.Category
	Article        *ent.Article
}

func setUpRecords(t *testing.T, entClient *ent.Client) records {
	t.Helper()

	categoryDocker := fixtureCategory.CreateCategory(t, entClient, fixtureCategory.CreateCategoryInput{Name: "Docker"})
	categoryGo := fixtureCategory.CreateCategory(t, entClient, fixtureCategory.CreateCategoryInput{Name: "Go"})

	article := fixtureArticle.CreateArticle(t, entClient, fixtureArticle.CreateArticleInput{
		Title:        "old title",
		BodyMarkdown: "old",
		BodyHTML:     "<p>old</p>\n",
		IsPublished:  false,
		Categories:   []*ent.Category{categoryDocker},
	})

	return records{
		CategoryDocker: categoryDocker,
		CategoryGo:     categoryGo,
		Article:        article,
	}
}

func callEndpoint(t *testing.T, server *httptest.Server, articleID uint32, body map[string]any) inertiaAction.TestActionResponse {
	t.Helper()

	path := "/admin/article/" + articleIDString(articleID)
	return inertiaAction.Send(t, server, inertiaAction.TestActionRequest{
		Method:       http.MethodPost,
		Path:         path,
		Body:         body,
		UseBasicAuth: true,
		Referer:      path,
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

func getArticleCategoryIDs(article *ent.Article) []uint32 {
	categoryIDs := make([]uint32, 0, len(article.Edges.Categories))
	for _, category := range article.Edges.Categories {
		categoryIDs = append(categoryIDs, category.ID)
	}
	sort.Slice(categoryIDs, func(i, j int) bool {
		return categoryIDs[i] < categoryIDs[j]
	})
	return categoryIDs
}

func articleIDString(articleID uint32) string {
	return strconv.FormatUint(uint64(articleID), 10)
}

func categoryIDString(categoryID uint32) string {
	return strconv.FormatUint(uint64(categoryID), 10)
}
