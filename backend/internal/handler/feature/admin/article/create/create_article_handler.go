package create

import (
	"net/http"

	"blog/internal/handler/handlererror"
	"blog/internal/handler/handlerresult"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Handle(r *http.Request) (handlerresult.HandlerResult, *handlererror.DisplayableError) {
	return handlerresult.Page("admin/CreateArticle", nil), nil
}
