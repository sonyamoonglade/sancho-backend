package validation

import (
	"testing"

	"github.com/sonyamoonglade/sancho-backend/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestValidatePayMethod(t *testing.T) {
	t.Run("should return ok", func(t *testing.T) {
		method := domain.PayOnline
		ok, msg := ValidatePayMethod(method)
		require.True(t, ok)
		require.Zero(t, msg)
	})

	t.Run("should return ok", func(t *testing.T) {
		method := domain.PayOnPickup
		ok, msg := ValidatePayMethod(method)
		require.True(t, ok)
		require.Zero(t, msg)
	})

	t.Run("should return err because no such method exist", func(t *testing.T) {
		method := domain.Pay{P: "some-random-bullshit"}
		ok, msg := ValidatePayMethod(method)
		require.False(t, ok)
		require.NotZero(t, msg)
		require.Equal(t, invalidPayMethod, msg)
	})
}
