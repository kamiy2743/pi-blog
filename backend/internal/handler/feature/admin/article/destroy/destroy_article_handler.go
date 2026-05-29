package destroy

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

func (h *Handler) Handle(r *http.Request) (handlerresult.ActionResult, error) {
	articleID, parseErr := article.ParseArticleID(r.PathValue("articleId"))
	if parseErr != nil {
		return handlerresult.ActionResult{}, &handlererror.DisplayableError{
			StatusCode: 404,
			Message:    "記事が見つかりません。",
			Err:        parseErr,
		}
	}

	if err := h.usecase.run(r.Context(), articleID); err != nil {
		return handlerresult.ActionResult{}, err
	}

	return handlerresult.Redirect("/admin", "記事を削除しました。"), nil
}
