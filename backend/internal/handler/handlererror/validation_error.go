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
