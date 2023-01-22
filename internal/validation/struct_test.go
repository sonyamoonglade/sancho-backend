package validation

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStructValidation(t *testing.T) {
	var (
		ok  bool
		msg string
	)
	type A struct {
		Name string `validate:"required"`
	}

	a := A{
		Name: "john",
	}
	ok, msg = ValidateStruct(a)
	require.True(t, ok)
	require.Zero(t, msg)

	a1 := A{}
	ok, msg = ValidateStruct(a1)
	require.False(t, ok)
	expected := `field "name" is missing in request body`
	require.Equal(t, expected, msg)
}
