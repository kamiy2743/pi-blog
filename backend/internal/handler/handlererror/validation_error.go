package handlererror

import (
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
