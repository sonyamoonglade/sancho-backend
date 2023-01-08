package fixtures

import (
	f "github.com/brianvoe/gofakeit/v6"
	"github.com/sonyamoonglade/sancho-backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetCategories(n int) []domain.Category {
	var categories []domain.Category
	for i := 0; i < n; i++ {
		categories = append(categories, generateCategory())
	}
	return categories
}

func generateCategory() domain.Category {
	return domain.Category{
		CategoryID: primitive.NewObjectID(),
		Rank:       int32(f.Number(1, 10)),
		Name:       f.BeerName(),
	}
}
