package service

import (
	"context"

	"github.com/sonyamoonglade/sancho-backend/internal/domain"
	"github.com/sonyamoonglade/sancho-backend/internal/services/dto"
)

type Product interface {
	GetByID(ctx context.Context, productID string) (domain.Product, error)
	GetAll(ctx context.Context) ([]domain.Product, error)
	GetAllCategories(ctx context.Context, sorted bool) ([]domain.Category, error)
	Create(ctx context.Context, dto dto.CreateProductDTO) (string, error)
	Delete(ctx context.Context, productID string) error
	Update(ctx context.Context, dto dto.UpdateProductDTO) error
	Approve(ctx context.Context, productID string) error
	Disapprove(ctx context.Context, productID string) error
}

type Auth interface {
	GetAdminByRefreshToken(ctx context.Context, token string) (domain.Admin, error)
	//GetCustomerByRefreshToken(token string)
	//GetWorkerByRefreshToken(token string)
}
