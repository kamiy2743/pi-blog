package edit

import (
	"net/http"

	"blog/internal/handler/inertia"

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
	result, err := h.usecase.run(r.Context())
	if err != nil {
		inertia.RenderError(w, r, h.inertia, *err)
		return
	}

	inertia.Render(w, r, h.inertia, "admin/EditCategory", format(result))
}
