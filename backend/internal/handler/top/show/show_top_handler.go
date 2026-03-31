package show

import (
	"net/http"

	"blog/internal/domain/article"
	"blog/internal/domain/category"

	"github.com/romsar/gonertia/v2"
)

func Handle(
	i *gonertia.Inertia,
	articleRepository article.ArticleRepository,
	categoryRepository category.CategoryRepository,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := Run(r.Context(), articleRepository, categoryRepository)
		if err != nil {
			http.Error(w, "記事取得エラー", http.StatusInternalServerError)
			return
		}
		if err := i.Render(w, r, "ShowTop", Format(result)); err != nil {
			http.Error(w, "描画エラー", http.StatusInternalServerError)
		}
	}
}
