package create

import (
	"net/http"

	"blog/internal/handler/handlerresult"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Handle(r *http.Request) (handlerresult.HandlerResult, error) {
	return handlerresult.Page("admin/CreateArticle", nil), nil
}
