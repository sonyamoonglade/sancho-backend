package appErrors

import "strings"

type AppError struct {
	stack    []string
	Original error
}

func (ap AppError) Error() string {
	return ap.Original.Error()
}

func (ap AppError) PrintStack() string {
	return strings.Join(ap.stack, ";")
}

func (ap AppError) OriginalError() error {
	return ap.Original
}

func WithContext(context string, err error) AppError {
	ap, ok := err.(*AppError)
	if !ok {
		return AppError{
			stack:    []string{context},
			Original: err,
		}
	}
	return AppError{
		stack:    append([]string{context}, ap.stack...),
		Original: ap.Original,
	}
}
