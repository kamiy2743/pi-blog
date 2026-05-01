package store

import "blog/internal/handler/handlererror"

func formatValidationError(validationError *handlererror.ValidationError) *handlererror.ValidationError {
	return handlererror.RemapValidationError(validationError, func(field string) string {
		return "create." + field
	})
}
