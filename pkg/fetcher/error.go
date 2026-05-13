package fetcher

import "fmt"

type validationError struct {
	msg string
}

func (e *validationError) Error() string {
	return e.msg
}

func ValidationError(pattern string, args ...any) error {
	return &validationError{msg: fmt.Sprintf(pattern, args...)}
}
