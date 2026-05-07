package destroy

import (
	"net/http"

	"blog/internal/domain/category"
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

func (h *Handler) Handle(r *http.Request) (handlerresult.ActionResult, error) {
	categoryID, parseErr := category.ParseCategoryID(r.PathValue("categoryId"))
	if parseErr != nil {
		return handlerresult.ActionResult{}, &handlererror.DisplayableError{
			StatusCode: 404,
			Message:    "カテゴリが見つかりません。",
			Err:        parseErr,
		}
	}

	if err := h.usecase.run(r.Context(), categoryID); err != nil {
		return handlerresult.ActionResult{}, err
	}

	return handlerresult.Redirect("/admin/category", "カテゴリを削除しました。"), nil
}
