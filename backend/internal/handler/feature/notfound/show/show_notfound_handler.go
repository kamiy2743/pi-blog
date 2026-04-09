package show

import (
	"net/http"

	"blog/internal/handler/inertia"

	"github.com/romsar/gonertia/v2"
)

type Handler struct {
	inertia *gonertia.Inertia
}

func NewHandler(i *gonertia.Inertia) *Handler {
	return &Handler{
		inertia: i,
	}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	inertia.RenderWithStatus(w, r, h.inertia, http.StatusNotFound, "NotFound", nil)
}
