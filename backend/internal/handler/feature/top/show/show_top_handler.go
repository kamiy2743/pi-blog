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

func (h *Handler) Handle(r *http.Request) (handlerresult.HandlerResult, error) {
	if r.URL.Path != "/" {
		return nil, &handlererror.DisplayableError{
			StatusCode:  http.StatusNotFound,
			Message:     "ページが見つかりません。",
			Description: "URL が変わったか、公開が終了した可能性があります。",
		}
	}

	result, err := h.usecase.run(r.Context())
	if err != nil {
		return nil, err
	}

	return handlerresult.Page("top/ShowTop", format(result)), nil
}
