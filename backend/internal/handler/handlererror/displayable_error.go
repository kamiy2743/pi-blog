package handlererror

import (
	"errors"
	"fmt"
)

type DisplayableError struct {
	StatusCode int
	Message    string
	Err        error
}

func (e DisplayableError) Error() string {
	return fmt.Sprintf(
		"statusCode=%d message=%q err=%v",
		e.StatusCode,
		e.Message,
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
