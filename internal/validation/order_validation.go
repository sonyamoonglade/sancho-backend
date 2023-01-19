package validation

import "github.com/sonyamoonglade/sancho-backend/internal/domain"

const (
	invalidPayMethod = "invalid pay method"
)

func ValidatePayType(p domain.Pay) (ok bool, msg string) {
	if p != domain.PayOnline || p != domain.PayOnPickup {
		return false, invalidPayMethod
	}

	return true, ""
}
