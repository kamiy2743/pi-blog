package search

import (
	"context"
	"net/http"

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
	prepareResult := inertia.PrepareInput(w, r, h.inertia, toInput)
	if !prepareResult.OK {
		return
	}
	input := prepareResult.Input

	inertia.Render(w, prepareResult.Request, h.inertia, "article/ShowArticleList", gonertia.Props{
		"initial": func(ctx context.Context) (any, error) {
			result, err := h.usecase.runInitial(ctx)
			if err != nil {
				return nil, err
			}
			return formatInitial(result), nil
		},
		"partialSearch": func(ctx context.Context) (any, error) {
			if prepareResult.HasValidationError {
				return formatPartialSearch(partialSearchResult{
					Title:       input.Title,
					CategoryIDs: input.CategoryIDs,
					Page:        1,
					TotalCount:  0,
					TotalPages:  1,
				}), nil
			}

			result, err := h.usecase.runPartialSearch(ctx, input)
			if err != nil {
				return nil, err
			}
			return formatPartialSearch(result), nil
		},
	})
}
