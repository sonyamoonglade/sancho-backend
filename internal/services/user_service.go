package service

import (
	"context"

	"github.com/sonyamoonglade/sancho-backend/internal/domain"
	"github.com/sonyamoonglade/sancho-backend/internal/services/dto"
	storage "github.com/sonyamoonglade/sancho-backend/internal/storages"
)

type Hasher interface {
	Hash(password string) string
}

type userService struct {
	passwordHasher Hasher
	userStorage    storage.User
}

func NewUserService(userStorage storage.User, hasher Hasher) User {
	return &userService{
		userStorage:    userStorage,
		passwordHasher: hasher,
	}
}

func (u userService) GetAdminByLogin(ctx context.Context, login string) (domain.Admin, error) {
	//TODO implement me
	panic("implement me")
}

func (u userService) GetAdminByRefreshToken(ctx context.Context, adminID, token string) (domain.Admin, error) {
	//TODO implement me
	panic("implement me")
}

func (u userService) SaveAdmin(ctx context.Context, admin domain.Admin) (string, error) {
	return u.userStorage.SaveAdmin(ctx, admin)
}

func (u userService) SaveSession(ctx context.Context, dto dto.SaveSessionDTO) error {
	return u.userStorage.SaveSession(ctx, dto)
}
