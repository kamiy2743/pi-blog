package handlererror

import (
	"errors"
	"fmt"
)

type DisplayableError struct {
	StatusCode  int
	Message     string
	Description string
	Err         error
}

func (e DisplayableError) Error() string {
	return fmt.Sprintf(
		"statusCode=%d message=%q description=%q err=%v",
		e.StatusCode,
		e.Message,
		e.Description,
		e.Err,
	)
}

func AsDisplayableError(err error) (*DisplayableError, bool) {
	var displayableError *DisplayableError
	if !errors.As(err, &displayableError) {
		return nil, false
	}
	return displayableError, true
}
