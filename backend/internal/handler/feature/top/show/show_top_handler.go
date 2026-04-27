package show

import (
	"net/http"

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
	if r.URL.Path != "/" {
		return handlerresult.PageResult{}, &handlererror.DisplayableError{
			StatusCode: http.StatusNotFound,
			Message:    "ページが見つかりません。",
		}
	}

	result, err := h.usecase.run(r.Context())
	if err != nil {
		return handlerresult.PageResult{}, err
	}

	return handlerresult.Page("top/ShowTop", format(result)), nil
}
