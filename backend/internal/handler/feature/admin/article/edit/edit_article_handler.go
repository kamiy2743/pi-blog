package edit

import (
	"net/http"

	"blog/internal/handler/handlerresult"

	"github.com/romsar/gonertia/v3"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Handle(r *http.Request) (handlerresult.PageResult, error) {
	return handlerresult.Page("admin/EditArticle", gonertia.Props{
		"articleId": r.PathValue("articleId"),
	}), nil
}
