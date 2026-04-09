package handlererror

type DisplayableError struct {
	StatusCode  int
	Message     string
	Description string
	Err         error
}
