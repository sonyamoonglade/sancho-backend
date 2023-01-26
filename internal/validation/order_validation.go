package validation

import (
	"github.com/sonyamoonglade/sancho-backend/internal/domain"
	"github.com/sonyamoonglade/sancho-backend/internal/handler/input"
)

const (
	invalidPayMethod = "invalid pay method"
	emptyCart        = "empty cart"
)

func ValidatePayMethod(p domain.Pay) (ok bool, msg string) {
	if p != domain.PayOnline && p != domain.PayOnPickup {
		return false, invalidPayMethod
	}

	return true, ""
}

func ValidateCart(cart []input.CartProductInput) (ok bool, msg string) {
	if len(cart) == 0 {
		return false, emptyCart
	}
	return true, ""
}
