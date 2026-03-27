package show

import (
	"net/http"

	inertia "github.com/romsar/gonertia/v2"
)

func Handle(i *inertia.Inertia) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := i.Render(w, r, "admin/ShowAdmin", nil); err != nil {
			http.Error(w, "描画エラー", http.StatusInternalServerError)
		}
	}
}
