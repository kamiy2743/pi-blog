package show

import (
	"net/http"

	"github.com/romsar/gonertia/v2"
)

type Handler struct {
	inertia *gonertia.Inertia
	usecase *Usecase
}

func NewHandler(i *gonertia.Inertia, u *Usecase) *Handler {
	return &Handler{
		inertia: i,
		usecase: u,
	}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	result, err := h.usecase.Run(r.Context())
	if err != nil {
		http.Error(w, "記事取得エラー", http.StatusInternalServerError)
		return
	}
	if err := h.inertia.Render(w, r, "ShowTop", Format(result)); err != nil {
		http.Error(w, "描画エラー", http.StatusInternalServerError)
	}
}
