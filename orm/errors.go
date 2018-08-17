package orm

import "fmt"

// NotFoundError is the error returned when an entity is not found when its expected to
type NotFoundError struct {
	entity string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s: Not Found", e.entity)
}
