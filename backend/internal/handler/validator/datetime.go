package validator

import (
	"blog/internal/domain"

	goValidator "github.com/go-playground/validator/v10"
)

func registerDatetime(validator *goValidator.Validate) {
	validator.RegisterValidation("datetime", func(fieldLevel goValidator.FieldLevel) bool {
		_, err := domain.ParseDatetime(fieldLevel.Field().String())
		return err == nil
	})

	validator.RegisterValidation("datetime_lt", func(fieldLevel goValidator.FieldLevel) bool {
		value := fieldLevel.Field().String()
		otherValue := fieldLevel.Parent().FieldByName(fieldLevel.Param()).String()
		if value == "" || otherValue == "" {
			return true
		}

		parsed, err := domain.ParseDatetime(value)
		if err != nil {
			return true
		}
		otherParsed, err := domain.ParseDatetime(otherValue)
		if err != nil {
			return true
		}

		return parsed.Before(otherParsed)
	})
}
