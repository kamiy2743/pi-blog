package handlererror

import (
	"errors"
	"fmt"
)

type ValidationError struct {
	Messages ValidationErrorMessages
}

type ValidationErrorMessages map[string]string

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error: %+v", e.Messages)
}

func (e *ValidationError) IsEmpty() bool {
	return e == nil || len(e.Messages) == 0
}

func (e *ValidationError) Merge(other *ValidationError) *ValidationError {
	if e == nil || e.IsEmpty() {
		return other
	}
	if other == nil || other.IsEmpty() {
		return e
	}

	for field, message := range other.Messages {
		e.Messages[field] = message
	}
	return e
}

func AsValidationError(err error) (*ValidationError, bool) {
	var validationError *ValidationError
	if !errors.As(err, &validationError) {
		return nil, false
	}
	if validationError == nil {
		return nil, false
	}
	return validationError, true
}
