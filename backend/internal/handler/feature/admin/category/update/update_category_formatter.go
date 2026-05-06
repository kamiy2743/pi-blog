package update

import (
	"strconv"

	"blog/internal/domain/category"
	"blog/internal/handler/handlererror"
)

func formatValidationError(
	validationError *handlererror.ValidationError,
	categoryID category.CategoryID,
) *handlererror.ValidationError {
	return handlererror.RemapValidationError(validationError, func(field string) string {
		return "update." + field + "." + strconv.FormatUint(uint64(categoryID), 10)
	})
}
