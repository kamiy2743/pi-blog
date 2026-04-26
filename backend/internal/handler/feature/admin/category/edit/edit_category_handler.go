package edit

import (
	"net/http"

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
	result, err := h.usecase.run(r.Context())
	if err != nil {
		return nil, err
	}

	return handlerresult.Page("admin/EditCategory", format(result)), nil
}
