package validation

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func ValidateStruct(v any) (ok bool, msg string) {
	if err := validate.Struct(v); err != nil {
		validationErr := err.(validator.ValidationErrors)
		firstErr := validationErr[0]
		msg = fmt.Sprintf("field %q is missing in request body", strings.ToLower(firstErr.Field()))
		return false, msg
	}
	return true, ""
}
