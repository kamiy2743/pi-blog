package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"blog/internal/model"

	inertia "github.com/romsar/gonertia/v2"
)

func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func ShowTop(i *inertia.Inertia) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := i.Render(w, r, "ShowTop", nil); err != nil {
			http.Error(w, "描画エラー", http.StatusInternalServerError)
		}
	}
}

func ShowNotFound(i *inertia.Inertia) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := renderWithStatus(w, http.StatusNotFound, func(target http.ResponseWriter) error {
			return i.Render(target, r, "NotFound", nil)
		}); err != nil {
			http.Error(w, "描画エラー", http.StatusInternalServerError)
		}
	}
}

func ShowArticleList(i *inertia.Inertia) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := i.Render(w, r, "article/ShowArticleList", nil); err != nil {
			http.Error(w, "描画エラー", http.StatusInternalServerError)
		}
	}
}

func ShowArticle(i *inertia.Inertia, articleID model.ArticleID) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		props := inertia.Props{
			"articleId": fmt.Sprintf("%s", articleID),
		}
		if err := i.Render(w, r, "article/ShowArticle", props); err != nil {
			http.Error(w, "描画エラー", http.StatusInternalServerError)
		}
	}
}

func ShowAdmin(i *inertia.Inertia) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := i.Render(w, r, "admin/ShowAdmin", nil); err != nil {
			http.Error(w, "描画エラー", http.StatusInternalServerError)
		}
	}
}

func CreateArticle(i *inertia.Inertia) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := i.Render(w, r, "admin/CreateArticle", nil); err != nil {
			http.Error(w, "描画エラー", http.StatusInternalServerError)
		}
	}
}

func StoreArticle(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func EditArticle(i *inertia.Inertia, articleID model.ArticleID) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		props := inertia.Props{
			"articleId": fmt.Sprintf("%s", articleID),
		}
		if err := i.Render(w, r, "admin/EditArticle", props); err != nil {
			http.Error(w, "描画エラー", http.StatusInternalServerError)
		}
	}
}

func UpdateArticle(w http.ResponseWriter, r *http.Request, articleID model.ArticleID) {
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func UpdatePublishSetting(w http.ResponseWriter, r *http.Request, articleID model.ArticleID) {
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func renderWithStatus(
	w http.ResponseWriter,
	statusCode int,
	render func(target http.ResponseWriter) error,
) error {
	recorder := httptest.NewRecorder()
	if err := render(recorder); err != nil {
		return err
	}

	for key, values := range recorder.Header() {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(statusCode)
	_, _ = w.Write(recorder.Body.Bytes())
	return nil
}
