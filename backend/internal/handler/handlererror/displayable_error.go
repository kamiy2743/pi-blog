package handlererror

import "fmt"

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
