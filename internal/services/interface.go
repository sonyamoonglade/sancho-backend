package service

import (
	"context"
	"time"

	"github.com/sonyamoonglade/sancho-backend/internal/domain"
	"github.com/sonyamoonglade/sancho-backend/internal/services/dto"
	"github.com/sonyamoonglade/sancho-backend/pkg/auth"
)

type Product interface {
	GetByID(ctx context.Context, productID string) (domain.Product, error)
	GetAll(ctx context.Context) ([]domain.Product, error)
	GetProductsByIDs(ctx context.Context, ids []string) ([]domain.Product, error)
	GetAllCategories(ctx context.Context, sorted bool) ([]domain.Category, error)
	Create(ctx context.Context, dto dto.CreateProductDTO) (string, error)
	Delete(ctx context.Context, productID string) error
	Update(ctx context.Context, dto dto.UpdateProductDTO) error
	Approve(ctx context.Context, productID string) error
	Disapprove(ctx context.Context, productID string) error
}

type Order interface {
	GetOrderByID(ctx context.Context, orderID string) (domain.Order, error)
	GetOrderByNanoIDAt(ctx context.Context, nanoID string, from, to time.Time) (domain.Order, error)
	GetLastOrderByCustomerID(ctx context.Context, customerID string) (domain.Order, error)

	CreateUserOrder(ctx context.Context, dto dto.CreateUserOrderDTO) (string, error)
	CreateWorkerOrder(ctx context.Context, orderDTO dto.CreateWorkerOrderDTO) (string, error)

	CalculateDiscountedAmount(amount int64, discountPercent float64) int64
	CalculateCartAmount(ctx context.Context, cart []dto.CartProductDTO) (int64, []domain.CartProduct, error)
}

type User interface {
	GetAdminByLogin(ctx context.Context, login string) (domain.Admin, error)
	GetAdminByRefreshToken(ctx context.Context, adminID, token string) (domain.Admin, error)
	GetCustomerByPhoneNumber(ctx context.Context, phoneNumber string) (domain.Customer, error)

	SaveAdmin(ctx context.Context, admin domain.Admin) (string, error)
	SaveCustomer(ctx context.Context, customer domain.Customer) (string, error)
	SaveSession(ctx context.Context, dto dto.SaveSessionDTO) error
}

type Auth interface {
	RegisterCustomer(ctx context.Context, dto dto.RegisterCustomerDTO) (string, error)
	RegisterAdmin(ctx context.Context, dto dto.RegisterAdminDTO) (string, error)
	LoginAdmin(ctx context.Context, dto dto.LoginAdminDTO) (auth.Pair, error)
	RefreshAdminToken(ctx context.Context, adminID, token string) (auth.Pair, error)
}
