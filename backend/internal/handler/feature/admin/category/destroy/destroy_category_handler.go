package destroy

import (
	"net/http"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/admin/category", http.StatusSeeOther)
}
