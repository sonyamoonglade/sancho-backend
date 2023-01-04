package validation

import (
	"github.com/sonyamoonglade/sancho-backend/internal/domain"
)

const (
	LiquidEmptyVolume         = "isLiquid but volume is 0"
	InsufficientEnergy        = "energy value can not be 0"
	InvalidNutrientsForLiquid = "liquid can not have nutrients"
)

func ValidateFeatures(f domain.Features) (ok bool, msg string) {
	if f.IsLiquid && f.Volume == 0 {
		return false, LiquidEmptyVolume
	}
	if f.EnergyValue == 0 {
		return false, InsufficientEnergy
	}
	if f.Nutrients != nil && f.IsLiquid {
		return false, InvalidNutrientsForLiquid
	}
	return true, ""
}
