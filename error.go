package httpclient

import "fmt"

type Error struct {
	Code    int
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("httpclient error. Code: %v detail: %s", e.Code, e.Message)
}
