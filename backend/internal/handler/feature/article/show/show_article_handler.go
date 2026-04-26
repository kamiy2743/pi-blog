package show

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
			StatusCode:  http.StatusNotFound,
			Message:     "ページが見つかりません。",
			Description: "URL が変わったか、公開が終了した可能性があります。",
			Err:         parseErr,
		}
	}

	result, usecaseErr := h.usecase.run(r.Context(), articleID)
	if usecaseErr != nil {
		return handlerresult.PageResult{}, usecaseErr
	}

	return handlerresult.Page("article/ShowArticle", format(result)), nil
}
