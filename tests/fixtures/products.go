package fixtures

import (
	f "github.com/brianvoe/gofakeit/v6"
	"github.com/sonyamoonglade/sancho-backend/internal/domain"
)

func GetLiquidFeatures() domain.Features {
	return domain.Features{
		IsLiquid:    true,
		Weight:      0,
		Volume:      int32(f.IntRange(100, 300)),
		EnergyValue: 0,
		Nutrients:   nil,
	}
}

func GetNonLiquidFeatures() domain.Features {
	return domain.Features{
		IsLiquid:    false,
		Weight:      int32(f.IntRange(100, 200)),
		Volume:      0,
		EnergyValue: int32(f.IntRange(200, 500)),
		Nutrients: &domain.Nutrients{
			Carbs:    f.Int32(),
			Proteins: f.Int32(),
			Fats:     f.Int32(),
		},
	}
}
