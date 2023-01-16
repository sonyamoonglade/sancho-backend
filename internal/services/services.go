package service

import (
	storage "github.com/sonyamoonglade/sancho-backend/internal/storages"
	"github.com/sonyamoonglade/sancho-backend/pkg/auth"
)

type Services struct {
	Product Product
	Auth    Auth
	User    User
}

type Deps struct {
	Storages      *storage.Storages
	TokenProvider auth.TokenProvider
	Hasher        Hasher
	TTLStrategy   TTLStrategy
}

func NewServices(deps Deps) *Services {
	stg := deps.Storages
	userService := NewUserService(stg.User, deps.Hasher)
	return &Services{
		Product: NewProductService(stg.Product),
		User:    userService,
		Auth:    NewAuthService(userService, deps.TokenProvider, deps.Hasher, deps.TTLStrategy),
	}
}
