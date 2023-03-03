package tests

import (
	"time"

	f "github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"github.com/sonyamoonglade/sancho-backend/internal/domain"
	service "github.com/sonyamoonglade/sancho-backend/internal/services"
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
			Name:        f.LoremIpsumSentence(4),
			TranslateRU: f.LoremIpsumSentence(5),
			Description: f.LoremIpsumSentence(10),
			ImageURL:    StringPtr(f.ImageURL(200, 200)),
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
			Name:        f.LoremIpsumSentence(2),
			TranslateRU: f.LoremIpsumSentence(5),
			Description: f.LoremIpsumSentence(10),
			ImageURL:    StringPtr(f.ImageURL(200, 200)),
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

	customer = domain.Customer{
		UserID:      primitive.NewObjectID(),
		PhoneNumber: "+79128557826",
		DeliveryAddress: &domain.UserDeliveryAddress{
			Address:   "Смирнова 20а",
			Entrance:  2,
			Floor:     9,
			Apartment: 73,
		},
		Role: domain.RoleCustomer,
		Name: StringPtr("Филипп"),
		Session: &domain.Session{
			RefreshToken: uuid.NewString(),
			ExpiresAt:    time.Now().UTC().Add(time.Hour * 24),
		},
	}

	worker = domain.Worker{
		UserID:   primitive.NewObjectID(),
		Name:     "Георгий",
		Login:    "jqkweixuch",
		Password: "kasdsjd*&1231mz",
		Role:     domain.RoleWorker,
		Session: domain.Session{
			RefreshToken: uuid.NewString(),
			ExpiresAt:    time.Now().UTC().Add(time.Hour * 24),
		},
	}

	meta = domain.BusinessMeta{
		DeliveryPunishmentThreshold: 500,
		DeliveryPunishmentValue:     100,
	}

	ttlStrategy = service.TTLStrategy{
		AccessTokenTTLs: map[domain.Role]time.Duration{
			domain.RoleAdmin:    time.Second * 60,
			domain.RoleWorker:   time.Second * 60,
			domain.RoleCustomer: time.Second * 3600,
		},
		RefreshTokenTTL: map[domain.Role]time.Duration{
			domain.RoleAdmin:  time.Hour * 1,
			domain.RoleWorker: time.Hour * 18,
			// ~ 1 year
			domain.RoleCustomer: time.Hour * 30 * 12,
		},
	}
)
