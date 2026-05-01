package validator

import (
	"strings"

	goValidator "github.com/go-playground/validator/v10"
)

func registerNotBlank(validator *goValidator.Validate) {
	validator.RegisterValidation("notblank", func(fieldLevel goValidator.FieldLevel) bool {
		value := fieldLevel.Field().String()
		return strings.TrimSpace(value) != ""
	})
}
