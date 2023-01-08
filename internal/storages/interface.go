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
	GetCategoryByName(ctx context.Context, categoryName string) (domain.Category, error)
	GetAllCategories(ctx context.Context, sorted bool) ([]domain.Category, error)
	Create(ctx context.Context, product domain.Product) (primitive.ObjectID, error)
	Delete(ctx context.Context, productID string) error
	Update(ctx context.Context, productID string, dto dto.UpdateProductDTO) error
	Approve(ctx context.Context, productID string) error
	Disapprove(ctx context.Context, productID string) error
	ChangeImageURL(ctx context.Context, productID string, imageURL string) error
}
