package fixtures

import (
	f "github.com/brianvoe/gofakeit/v6"
	"github.com/sonyamoonglade/sancho-backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetProducts(n int) []interface{} {
	var products []interface{}
	for i := 0; i < n; i++ {
		products = append(products, generateProduct())
	}
	return products
}

func GetProduct() domain.Product {
	return generateProduct().(domain.Product)
}

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

func generateProduct() interface{} {
	return domain.Product{
		ProductID:   primitive.NewObjectID(),
		Name:        f.Name(),
		TranslateRU: f.Name(),
		Description: f.LoremIpsumSentence(20),
		ImageURL:    f.ImageURL(100, 200),
		Price:       f.Int64(),
		IsApproved:  false,
		Category:    generateCategory(),
		Features: domain.Features{
			Weight:      f.Int32(),
			EnergyValue: f.Int32(),
			Volume:      f.Int32(),
			Nutrients: &domain.Nutrients{
				Carbs:    f.Int32(),
				Proteins: f.Int32(),
				Fats:     f.Int32(),
			},
		},
	}
}
