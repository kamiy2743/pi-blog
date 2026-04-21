package handlererror

import (
	"fmt"
)

type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf(
		"field=%q message=%q",
		e.Field,
		e.Message,
	)
}

func ValidationErrorsToMap(validationErrors []ValidationError) map[string]string {
	errorsMap := make(map[string]string, len(validationErrors))
	for _, validationError := range validationErrors {
		errorsMap[validationError.Field] = validationError.Message
	}
	return errorsMap
}
