package utils

import "fmt"

// NameAlreadyTaken is a specific type of error used when
// a user tries to persist an entity under the same name
// of an already existing one.
type NameAlreadyTaken struct {
	Err error
}

func (e *NameAlreadyTaken) Error() string {
	return fmt.Sprintf("%d", e.Err)
}
