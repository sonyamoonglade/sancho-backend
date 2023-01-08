package tests

import (
	f "github.com/brianvoe/gofakeit/v6"
	"github.com/sonyamoonglade/sancho-backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	categoryPizza = domain.Category{
		CategoryID: primitive.NewObjectID(),
		Rank:       1,
		Name:       "Пицца",
	}
	categoryDrinks = domain.Category{
		CategoryID: primitive.NewObjectID(),
		Rank:       2,
		Name:       "Напитки",
	}

	categories = []interface{}{categoryDrinks, categoryPizza}

	products = []interface{}{
		domain.Product{
			ProductID:   primitive.NewObjectID(),
			Name:        f.BeerName(),
			TranslateRU: f.BeerName(),
			Description: f.LoremIpsumSentence(10),
			ImageURL:    f.ImageURL(200, 200),
			IsApproved:  false,
			Price:       int64(f.IntRange(200, 500)),
			Category:    categoryPizza,
			Features: domain.Features{
				IsLiquid:    false,
				Weight:      300,
				Volume:      0,
				EnergyValue: 250,
				Nutrients: &domain.Nutrients{
					Carbs:    35,
					Proteins: 22,
					Fats:     19,
				},
			},
		},
		domain.Product{
			ProductID:   primitive.NewObjectID(),
			Name:        f.BeerName(),
			TranslateRU: f.BeerName(),
			Description: f.LoremIpsumSentence(10),
			ImageURL:    f.ImageURL(200, 200),
			IsApproved:  true,
			Price:       int64(f.IntRange(50, 100)),
			Category:    categoryDrinks,
			Features: domain.Features{
				IsLiquid:    true,
				Weight:      200,
				Volume:      200,
				EnergyValue: 50,
				Nutrients:   nil,
			},
		},
	}
)
