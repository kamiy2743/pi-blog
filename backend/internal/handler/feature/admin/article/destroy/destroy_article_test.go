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
	domainArticle "blog/internal/domain/article"
	"blog/internal/test"
	fixtureArticle "blog/internal/test/fixture/article"
	fixtureCategory "blog/internal/test/fixture/category"
	fixtureComment "blog/internal/test/fixture/comment"
	"blog/internal/test/helper"
	inertiaAction "blog/internal/test/helper/inertia/action"
	stubArticle "blog/internal/test/stub/article"
)

func Test記事を削除できる(t *testing.T) {
	initResult := test.Init(t)
	article := setUpRecords(t, initResult.EntClient)

	res := callEndpoint(t, initResult.Server, article.ID)
	res.AssertRedirectTo(t, "/admin")

	assertArticleCount(t, initResult.EntClient, 0)
	assertCategoryCount(t, initResult.EntClient, 1)
	assertCommentCount(t, initResult.EntClient, 0)
}

func Test記事が存在しない場合はエラー(t *testing.T) {
	initResult := test.Init(t)

	res := callEndpoint(t, initResult.Server, 1)
	res.AssertRedirectTo(t, "/admin/article/1")
	res.AssertFlashError(t, initResult.Server, initResult.SessionManager, "記事が見つかりません。")

	assertArticleCount(t, initResult.EntClient, 0)
}

func Test記事IDが不正な場合はエラー(t *testing.T) {
	initResult := test.Init(t)

	res := inertiaAction.Send(t, initResult.Server, inertiaAction.TestActionRequest{
		Method:       http.MethodPost,
		Path:         "/admin/article/invalid/delete",
		Body:         map[string]any{},
		UseBasicAuth: true,
		Referer:      "/admin/article/invalid",
	})
	res.AssertRedirectTo(t, "/admin/article/invalid")
	res.AssertFlashError(t, initResult.Server, initResult.SessionManager, "記事が見つかりません。")
}

func Test記事の検索に失敗した場合はエラー(t *testing.T) {
	initResult := test.Init(t, &di.ContainerOptions{
		ArticleRepository: stubArticle.ArticleRepositoryStub{
			SearchFunc: func(ctx context.Context, criteria domainArticle.SearchArticleCriteria) ([]domainArticle.Article, error) {
				return nil, errors.New("test")
			},
		},
	})

	res := callEndpoint(t, initResult.Server, 1)
	res.AssertRedirectTo(t, "/admin/article/1")
	res.AssertFlashError(t, initResult.Server, initResult.SessionManager, "記事の削除に失敗しました。")
}

func Test記事の削除に失敗した場合はエラー(t *testing.T) {
	initResult := test.Init(t, &di.ContainerOptions{
		ArticleRepository: stubArticle.ArticleRepositoryStub{
			SearchFunc: func(ctx context.Context, criteria domainArticle.SearchArticleCriteria) ([]domainArticle.Article, error) {
				return []domainArticle.Article{
					{ID: 1, Title: "Go", BodyMarkdown: "本文", BodyHTML: "<p>本文</p>\n", UpdatedAt: helper.Time(t, "2026-01-01 00:00")},
				}, nil
			},
			DeleteFunc: func(ctx context.Context, article domainArticle.Article) error {
				return errors.New("test")
			},
		},
	})

	res := callEndpoint(t, initResult.Server, 1)
	res.AssertRedirectTo(t, "/admin/article/1")
	res.AssertFlashError(t, initResult.Server, initResult.SessionManager, "記事の削除に失敗しました。")
}

func callEndpoint(t *testing.T, server *httptest.Server, articleID uint32) inertiaAction.TestActionResponse {
	t.Helper()

	path := "/admin/article/" + strconv.FormatUint(uint64(articleID), 10)
	return inertiaAction.Send(t, server, inertiaAction.TestActionRequest{
		Method:       http.MethodPost,
		Path:         path + "/delete",
		Body:         map[string]any{},
		UseBasicAuth: true,
		Referer:      path,
	})
}

func setUpRecords(t *testing.T, entClient *ent.Client) *ent.Article {
	t.Helper()

	category := fixtureCategory.CreateCategory(t, entClient, fixtureCategory.CreateCategoryInput{Name: "Go"})
	article := fixtureArticle.CreateArticle(t, entClient, fixtureArticle.CreateArticleInput{
		Title:        "Go",
		BodyMarkdown: "本文",
		Categories:   []*ent.Category{category},
	})
	fixtureComment.CreateComment(t, entClient, fixtureComment.CreateCommentInput{
		Article:    article,
		AuthorName: "tester",
		Body:       "コメント",
	})
	return article
}

func assertArticleCount(t *testing.T, entClient *ent.Client, expected int) {
	t.Helper()

	count, err := entClient.Article.Query().Count(context.Background())
	if err != nil {
		t.Fatalf("記事件数の取得に失敗: %v", err)
	}
	helper.AssertEqual(t, expected, count, "記事件数が不正です")
}

func assertCategoryCount(t *testing.T, entClient *ent.Client, expected int) {
	t.Helper()

	count, err := entClient.Category.Query().Count(context.Background())
	if err != nil {
		t.Fatalf("カテゴリ件数の取得に失敗: %v", err)
	}
	helper.AssertEqual(t, expected, count, "カテゴリ件数が不正です")
}

func assertCommentCount(t *testing.T, entClient *ent.Client, expected int) {
	t.Helper()

	count, err := entClient.Comment.Query().Count(context.Background())
	if err != nil {
		t.Fatalf("コメント件数の取得に失敗: %v", err)
	}
	helper.AssertEqual(t, expected, count, "コメント件数が不正です")
}
