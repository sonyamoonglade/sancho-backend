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
		field := firstErr.Field()
		decamelizedField := decamelizeFirstLetter(field)
		// Decamelize first letter
		msg = fmt.Sprintf("field %q is missing in request body", decamelizedField)
		return false, msg
	}
	return true, ""
}

func decamelizeFirstLetter(word string) string {
	byteArr := []byte(word)
	firstLetterLower := strings.ToLower(string(byteArr[0]))
	return firstLetterLower + word[1:]
}
