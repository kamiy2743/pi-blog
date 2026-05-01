package store

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

func (h *Handler) Handle(r *http.Request) (handlerresult.ActionResult, error) {
	input, validationError := toInput(r)
	if validationError != nil {
		return handlerresult.ActionResult{}, formatValidationError(validationError)
	}

	if err := h.usecase.run(r.Context(), input); err != nil {
		return handlerresult.ActionResult{}, err
	}

	return handlerresult.Redirect("/admin/category", "カテゴリを作成しました。"), nil
}
