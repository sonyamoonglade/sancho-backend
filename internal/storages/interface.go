package storage

import (
	"context"

	"github.com/sonyamoonglade/sancho-backend/internal/domain"
	"github.com/sonyamoonglade/sancho-backend/internal/services/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product interface {
	GetByID(ctx context.Context, productID string) (domain.Product, error)
	GetAll(ctx context.Context) ([]domain.Product, error)
	Create(ctx context.Context, product domain.Product) (primitive.ObjectID, error)
	Update(ctx context.Context, dto dto.UpdateProductDTO) error
	Delete(ctx context.Context, productID string) error
	Approve(ctx context.Context, productID string) error
	Disapprove(ctx context.Context, productID string) error
	GetCategoryByName(ctx context.Context, categoryName string) (domain.Category, error)
	GetAllCategories(ctx context.Context, sorted bool) ([]domain.Category, error)
}
