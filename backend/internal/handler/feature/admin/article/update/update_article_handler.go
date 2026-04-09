package update

import (
	"net/http"

	"blog/internal/domain/article"
)

func Handle(w http.ResponseWriter, r *http.Request, articleID article.ArticleID) {
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
