package utils

import "fmt"

type UsernameTaken struct {
	Err error
}

func (e *UsernameTaken) Error() string {
	return fmt.Sprintf("%d", e.Err)
}
