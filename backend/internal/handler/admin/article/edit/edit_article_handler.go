package edit

import (
	"fmt"
	"net/http"

	"blog/internal/domain/article"

	inertia "github.com/romsar/gonertia/v2"
)

func Handle(i *inertia.Inertia, articleID article.ArticleID) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		props := inertia.Props{
			"articleId": fmt.Sprintf("%s", articleID),
		}
		if err := i.Render(w, r, "admin/EditArticle", props); err != nil {
			http.Error(w, "描画エラー", http.StatusInternalServerError)
		}
	}
}
