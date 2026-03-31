package edit

import (
	"net/http"

	"blog/internal/domain/article"

	"github.com/romsar/gonertia/v2"
)

func Handle(i *gonertia.Inertia, articleID article.ArticleID) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		props := gonertia.Props{
			"articleId": articleID,
		}
		if err := i.Render(w, r, "admin/EditArticle", props); err != nil {
			http.Error(w, "描画エラー", http.StatusInternalServerError)
		}
	}
}
