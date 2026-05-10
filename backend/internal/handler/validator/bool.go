package validator

import goValidator "github.com/go-playground/validator/v10"

func registerBool(validator *goValidator.Validate) {
	validator.RegisterValidation("bool", func(fieldLevel goValidator.FieldLevel) bool {
		value := fieldLevel.Field().String()
		return value == "true" || value == "false"
	})
}
