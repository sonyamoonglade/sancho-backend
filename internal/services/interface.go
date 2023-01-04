package service

import (
	"context"

	"github.com/sonyamoonglade/sancho-backend/internal/domain"
	"github.com/sonyamoonglade/sancho-backend/internal/services/dto"
)

type Product interface {
	GetAll(ctx context.Context) ([]domain.Product, error)
	GetAllCategories(ctx context.Context, sorted bool) ([]domain.Category, error)
	Create(ctx context.Context, dto dto.CreateProductDTO) error
	Delete(ctx context.Context, productID string) error
	Update(ctx context.Context, productID string, dto dto.UpdateProductDTO) error
	Approve(ctx context.Context, productID string) error
	ChangeImageURL(ctx context.Context, productID string, imageURL string) error
}
