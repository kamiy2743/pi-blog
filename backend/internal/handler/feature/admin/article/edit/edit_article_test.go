package edit_test

import (
	"context"
	"errors"
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
	inertiaPage "blog/internal/test/helper/inertia/page"
	stubArticle "blog/internal/test/stub/article"
	stubCategory "blog/internal/test/stub/category"

	"github.com/romsar/gonertia/v3"
)

func Test記事編集画面を表示できる(t *testing.T) {
	initResult := test.Init(t)
	records := setUpRecords(t, initResult.EntClient)

	res := callEndpoint(t, initResult.Server, records.Article.ID)

	res.AssertFullProps(t, "admin/EditArticle", gonertia.Props{
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
		"article": gonertia.Props{
			"id":             domainArticle.ArticleID(records.Article.ID),
			"title":          "old title",
			"bodyMarkdown":   "# old",
			"isPublished":    true,
			"publishStartAt": "2026-01-02T03:04:00Z",
			"publishEndAt":   "2026-01-03T04:05:00Z",
			"categoryIds":    []uint32{records.CategoryDocker.ID},
		},
	})
}

func Test存在しない記事なら404(t *testing.T) {
	initResult := test.Init(t)

	res := callEndpoint(t, initResult.Server, 999)

	res.AssertError(t, 404, "記事が見つかりません。")
}

func Testカテゴリの読み込みに失敗した場合は500(t *testing.T) {
	initResult := test.Init(t, &di.ContainerOptions{
		CategoryRepository: stubCategory.CategoryRepositoryStub{
			AllFunc: func(ctx context.Context, orderBy domainCategory.OrderBy) ([]domainCategory.Category, error) {
				return nil, errors.New("test")
			},
		},
	})

	res := callEndpoint(t, initResult.Server, 1)

	res.AssertError(t, 500, "カテゴリの読み込みに失敗しました。")
}

func Test記事の読み込みに失敗した場合は500(t *testing.T) {
	initResult := test.Init(t, &di.ContainerOptions{
		ArticleRepository: stubArticle.ArticleRepositoryStub{
			SearchFunc: func(ctx context.Context, criteria domainArticle.SearchArticleCriteria) ([]domainArticle.Article, error) {
				return nil, errors.New("test")
			},
		},
	})

	res := callEndpoint(t, initResult.Server, 1)

	res.AssertError(t, 500, "記事の読み込みに失敗しました。")
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
		Title:          "old title",
		BodyMarkdown:   "# old",
		BodyHTML:       "<h1>old</h1>\n",
		IsPublished:    true,
		PublishStartAt: helper.TimePtr(t, "2026-01-02 03:04"),
		PublishEndAt:   helper.TimePtr(t, "2026-01-03 04:05"),
		Categories:     []*ent.Category{categoryDocker},
	})

	return records{
		CategoryDocker: categoryDocker,
		CategoryGo:     categoryGo,
		Article:        article,
	}
}

func callEndpoint(t *testing.T, server *httptest.Server, articleID uint32) inertiaPage.TestPageResponse {
	t.Helper()

	return inertiaPage.Send(t, server, inertiaPage.TestPageRequest{
		Path:         "/admin/article/" + strconv.FormatUint(uint64(articleID), 10),
		UseBasicAuth: true,
	})
}
