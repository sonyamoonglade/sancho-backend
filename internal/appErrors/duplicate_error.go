package appErrors

import "net/http"

const (
	with          = " with "
	alreadyExists = " already exists"
)

// DuplicateError signals if an update or insert did duplicate update/insert.
type DuplicateError struct {
	msg  string
	code int
}

func (u DuplicateError) Error() string {
	return u.msg
}

func (u DuplicateError) Code() int {
	return u.code
}

// NewDuplicateError returns instance of UpdateError.
// entity is for example product, category, order etc...
// field is on what field duplicate update occurred.
// value is field update value.
func NewDuplicateError(entity, field, value string) error {
	return DuplicateError{
		msg:  entity + with + field + " " + value + alreadyExists,
		code: http.StatusConflict,
	}
}
