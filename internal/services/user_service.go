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
	adminID, err := u.userStorage.SaveAdmin(ctx, admin)
	if err != nil {
		return "", err
	}
	return adminID.Hex(), nil
}

func (u userService) GetCustomerByPhoneNumber(ctx context.Context, phoneNumber string) (domain.Customer, error) {
	return u.userStorage.GetCustomerByPhoneNumber(ctx, phoneNumber)
}

func (u userService) SaveCustomer(ctx context.Context, customer domain.Customer) (string, error) {
	customerID, err := u.userStorage.SaveCustomer(ctx, customer)
	if err != nil {
		return "", err
	}
	return customerID.Hex(), err
}

func (u userService) SaveSession(ctx context.Context, dto dto.SaveSessionDTO) error {
	return u.userStorage.SaveSession(ctx, dto)
}
