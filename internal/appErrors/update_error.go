package appErrors

import "net/http"

const (
	with          = " with "
	alreadyExists = " already exists"
)

// UpdateError signals if an update updated field and it became duplicate.
type UpdateError struct {
	msg  string
	code int
}

func (u UpdateError) Error() string {
	return u.msg
}

func (u UpdateError) Code() int {
	return u.code
}

// NewUpdateError returns instance of UpdateError.
// entity is for example product, category, order etc...
// field is on what field duplicate update occurred.
// value is field update value.
func NewUpdateError(entity, field, value string) error {
	return UpdateError{
		msg:  entity + with + field + " " + value + alreadyExists,
		code: http.StatusConflict,
	}
}
