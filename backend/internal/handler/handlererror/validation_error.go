package handlererror

import (
	"errors"
	"fmt"
)

type ValidationError struct {
	Messages map[string]string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error: %+v", e.Messages)
}

func (e *ValidationError) IsEmpty() bool {
	return e == nil || len(e.Messages) == 0
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

func RemapValidationError(
	validationError *ValidationError,
	remap func(field string) string,
) *ValidationError {
	if validationError == nil {
		return nil
	}

	remappedMessages := make(map[string]string, len(validationError.Messages))
	for field, message := range validationError.Messages {
		remappedMessages[remap(field)] = message
	}

	return &ValidationError{
		Messages: remappedMessages,
	}
}
