package validator

import (
	"errors"

	"blog/internal/handler/handlererror"

	goValidator "github.com/go-playground/validator/v10"
)

func Validate[T any](
	req T,
	toValidationErr func(field, tag string) *handlererror.ValidationError,
) []handlererror.ValidationError {
	err := goValidator.New().Struct(req)
	if err == nil {
		return nil
	}

	var validationErrors goValidator.ValidationErrors
	errors.As(err, &validationErrors)

	result := make([]handlererror.ValidationError, 0, len(validationErrors))
	for _, validationError := range validationErrors {
		err := toValidationErr(validationError.Field(), validationError.Tag())
		if err == nil {
			result = append(result, handlererror.ValidationError{
				Field:   validationError.Field(),
				Message: "入力内容が不正です。",
			})
		}
		result = append(result, *err)
	}

	return result
}
