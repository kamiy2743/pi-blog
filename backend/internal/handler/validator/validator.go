package validator

import (
	"errors"
	"reflect"
	"strings"

	"blog/internal/handler/handlererror"

	goValidator "github.com/go-playground/validator/v10"
)

func newValidator() *goValidator.Validate {
	validator := goValidator.New()
	validator.RegisterTagNameFunc(toLowerCamelCase)
	registerBool(validator)
	registerDatetime(validator)
	registerNotBlank(validator)
	return validator
}

func Validate[T any](
	req T,
	getValidationMessage func(field, tag string) string,
) *handlererror.ValidationError {
	err := newValidator().Struct(req)
	if err == nil {
		return nil
	}

	var validationErrors goValidator.ValidationErrors
	errors.As(err, &validationErrors)

	messages := make(handlererror.ValidationErrorMessages, len(validationErrors))
	for _, validationError := range validationErrors {
		field := validationError.Field()
		message := getValidationMessage(field, validationError.Tag())
		messages[field] = message
	}
	return &handlererror.ValidationError{Messages: messages}
}

func toLowerCamelCase(field reflect.StructField) string {
	name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
	if name == "" || name == "-" {
		return strings.ToLower(field.Name)
	}
	return name
}
