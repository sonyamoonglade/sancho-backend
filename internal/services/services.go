package service

import (
	"github.com/sonyamoonglade/sancho-backend/internal/domain"
	storage "github.com/sonyamoonglade/sancho-backend/internal/storages"
	"github.com/sonyamoonglade/sancho-backend/pkg/auth"
)

type Services struct {
	Product Product
	Auth    Auth
	User    User
	Order   Order
}

type Deps struct {
	Storages      *storage.Storages
	TokenProvider auth.TokenProvider
	MetaProvider  domain.MetaProvider
	Hasher        Hasher
	TTLStrategy   TTLStrategy
	OrderConfig   OrderConfig
}

func NewServices(deps Deps) *Services {
	stg := deps.Storages
	userService := NewUserService(stg.User, deps.Hasher)
	productService := NewProductService(stg.Product)
	return &Services{
		Product: productService,
		User:    userService,
		Auth:    NewAuthService(userService, deps.TokenProvider, deps.Hasher, deps.TTLStrategy),
		Order:   NewOrderService(stg.Order, productService, deps.OrderConfig, deps.MetaProvider),
	}
}
