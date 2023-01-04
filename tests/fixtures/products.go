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

func generateProduct() interface{} {
	return domain.Product{
		ProductID:   primitive.NewObjectID().String(),
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
