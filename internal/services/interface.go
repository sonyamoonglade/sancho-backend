package service

import (
	"context"

	"github.com/sonyamoonglade/sancho-backend/auth"
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
	RegisterAdmin(ctx context.Context, dto dto.RegisterAdminDTO) (string, error)
	LoginAdmin(ctx context.Context, dto dto.LoginAdminDTO) (auth.Pair, error)
	RefreshAdminToken(ctx context.Context, adminID, token string) (auth.Pair, error)
}

type User interface {
	SaveAdmin(ctx context.Context, admin domain.Admin) (string, error)
	GetAdminByLogin(ctx context.Context, login string) (domain.Admin, error)
	GetAdminByRefreshToken(ctx context.Context, adminID, token string) (domain.Admin, error)

	SaveSession(ctx context.Context, dto dto.SaveSessionDTO) error
}
