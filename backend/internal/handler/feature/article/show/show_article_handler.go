package show

import (
	"net/http"

	"blog/internal/domain/article"
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
	articleID, err := article.ParseArticleID(r.PathValue("articleId"))
	if err != nil {
		inertia.RenderNotFound(w, r, h.inertia)
		return
	}

	result, usecaseErr := h.usecase.run(r.Context(), articleID)
	if usecaseErr != nil {
		inertia.RenderError(w, r, h.inertia, *usecaseErr)
		return
	}

	inertia.Render(w, r, h.inertia, "article/ShowArticle", format(result))
}
