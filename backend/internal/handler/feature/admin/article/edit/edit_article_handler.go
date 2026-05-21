package edit

import (
	"net/http"

	"blog/internal/domain/article"
	"blog/internal/handler/handlererror"
	"blog/internal/handler/handlerresult"
)

type Handler struct {
	usecase *Usecase
}

func NewHandler(u *Usecase) *Handler {
	return &Handler{
		usecase: u,
	}
}

func (h *Handler) Handle(r *http.Request) (handlerresult.PageResult, error) {
	articleID, parseErr := article.ParseArticleID(r.PathValue("articleId"))
	if parseErr != nil {
		return handlerresult.PageResult{}, &handlererror.DisplayableError{
			StatusCode: 404,
			Message:    "記事が見つかりません。",
			Err:        parseErr,
		}
	}

	result, err := h.usecase.run(r.Context(), articleID)
	if err != nil {
		return handlerresult.PageResult{}, err
	}

	return handlerresult.Page("admin/EditArticle", format(result)), nil
}
