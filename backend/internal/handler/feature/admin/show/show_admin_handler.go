package show

import (
	"net/http"

	"blog/internal/handler/handlererror"
	"blog/internal/handler/handlerresult"
	"blog/internal/handler/inertia"

	"github.com/romsar/gonertia/v3"
)

const component = "admin/ShowAdmin"

type Handler struct {
	usecase *Usecase
}

func NewHandler(u *Usecase) *Handler {
	return &Handler{
		usecase: u,
	}
}

func (h *Handler) Handle(r *http.Request) (handlerresult.PageResult, error) {
	input, validationError := toInput(r)
	props := gonertia.Props{}

	if inertia.ShouldIncludeProp(r, component, "initial") {
		result, err := h.usecase.runInitial(r.Context())
		if err != nil {
			return handlerresult.PageResult{}, err
		}
		props["initial"] = formatInitial(result)
	}

	if inertia.ShouldIncludeProp(r, component, "partialSearch") {
		partialSearchProps, err := h.handlePartialSearch(r, input, validationError != nil)
		props["partialSearch"] = partialSearchProps
		if err != nil {
			return handlerresult.Page(component, props), err
		}
	}

	if validationError == nil {
		return handlerresult.Page(component, props), nil
	}
	return handlerresult.Page(component, props), validationError
}

func (h *Handler) handlePartialSearch(
	r *http.Request,
	input input,
	hasValidationError bool,
) (gonertia.Props, *handlererror.DisplayableError) {
	if hasValidationError {
		return formatPartialSearch(emptyPartialSearchResult(input)), nil
	}

	result, err := h.usecase.runPartialSearch(r.Context(), input)
	if err != nil {
		return formatPartialSearch(emptyPartialSearchResult(input)), err
	}
	return formatPartialSearch(result), nil
}

func emptyPartialSearchResult(input input) partialSearchResult {
	return partialSearchResult{
		Title:       input.Title,
		CategoryIDs: input.CategoryIDs,
		Page:        1,
		TotalCount:  0,
		TotalPages:  1,
	}
}
