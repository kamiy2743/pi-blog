package update

import (
	"net/http"

	"blog/internal/handler/handlerresult"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Handle(r *http.Request) (handlerresult.ActionResult, error) {
	return handlerresult.Redirect("/admin", "Success"), nil
}
