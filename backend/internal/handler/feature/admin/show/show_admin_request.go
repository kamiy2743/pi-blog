package show

import (
	"net/http"
	"strconv"

	"blog/internal/domain/category"
	"blog/internal/handler/handlererror"
	"blog/internal/handler/validator"
)

const perPage = 20

type request struct {
	Title       string `validate:"omitempty,max=255"`
	CategoryIDs []string
	Page        string
}

func toInput(r *http.Request) (input, []handlererror.ValidationError, *handlererror.DisplayableError) {
	req := request{
		Title:       r.URL.Query().Get("title"),
		CategoryIDs: r.URL.Query()["categoryId"],
		Page:        r.URL.Query().Get("page"),
	}
	validationErrs := validator.Validate(req, toValidationError)

	categoryIDs := parseCategoryIDs(req.CategoryIDs)
	page := parsePage(req.Page)

	return input{
		Title:       req.Title,
		CategoryIDs: categoryIDs,
		Page:        page,
		PerPage:     perPage,
	}, validationErrs, nil
}

func parseCategoryIDs(categoryIDStrs []string) []category.CategoryID {
	categoryIDs := make([]category.CategoryID, 0, len(categoryIDStrs))
	for _, categoryIDStr := range categoryIDStrs {
		categoryID, err := category.ParseCategoryID(categoryIDStr)
		if err != nil {
			continue
		}
		categoryIDs = append(categoryIDs, categoryID)
	}
	return categoryIDs
}

func parsePage(pageStr string) int {
	if pageStr == "" {
		return 1
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		return 1
	}
	return page
}

func toValidationError(field, tag string) *handlererror.ValidationError {
	switch field {
	case "title":
		if tag == "max" {
			return &handlererror.ValidationError{
				Field:   field,
				Message: "タイトルは255文字以下で入力してください。",
			}
		}
	}
	return nil
}
