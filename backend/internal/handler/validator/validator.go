package validator

import (
	"errors"
	"strings"

	"blog/internal/handler/handlererror"

	goValidator "github.com/go-playground/validator/v10"
)

func Validate[T any](
	req T,
	getValidationMessage func(field, tag string) string,
) *handlererror.ValidationError {
	err := goValidator.New().Struct(req)
	if err == nil {
		return nil
	}

	var validationErrors goValidator.ValidationErrors
	errors.As(err, &validationErrors)

	messages := make(map[string]string, len(validationErrors))
	for _, validationError := range validationErrors {
		field := strings.ToLower(validationError.Field())
		message := getValidationMessage(field, validationError.Tag())
		messages[field] = message
	}
	return &handlererror.ValidationError{Messages: messages}
}
