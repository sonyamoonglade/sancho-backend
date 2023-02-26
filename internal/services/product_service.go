package service

import (
	"context"

	"github.com/sonyamoonglade/sancho-backend/internal/appErrors"
	"github.com/sonyamoonglade/sancho-backend/internal/domain"
	"github.com/sonyamoonglade/sancho-backend/internal/services/dto"
	"github.com/sonyamoonglade/sancho-backend/internal/storages"
)

type productService struct {
	productStorage storage.Product
}

func NewProductService(productStorage storage.Product) Product {
	return &productService{productStorage: productStorage}
}

func (p productService) GetByID(ctx context.Context, productID string) (domain.Product, error) {
	product, err := p.productStorage.GetByID(ctx, productID)
	if err != nil {
		return domain.Product{}, err
	}
	return product, nil
}

func (p productService) GetAll(ctx context.Context) ([]domain.Product, error) {
	catalog, err := p.productStorage.GetAll(ctx)
	if err != nil {
		return nil, appErrors.WithContext("productStorage.GetAll", err)
	}
	return catalog, nil
}

func (p productService) GetProductsByIDs(ctx context.Context, ids []string) ([]domain.Product, error) {
	return p.productStorage.GetByIDs(ctx, ids)
}

func (p productService) GetAllCategories(ctx context.Context, sorted bool) ([]domain.Category, error) {
	categories, err := p.productStorage.GetAllCategories(ctx, sorted)
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (p productService) Create(ctx context.Context, dto dto.CreateProductDTO) (string, error) {
	category, err := p.productStorage.GetCategoryByName(ctx, dto.CategoryName)
	if err != nil {
		return "", err
	}

	product := dto.ToDomain()
	product.Category = category

	productID, err := p.productStorage.Save(ctx, product)
	if err != nil {
		return "", err
	}

	return productID.Hex(), nil
}

func (p productService) Delete(ctx context.Context, productID string) error {
	return p.productStorage.Delete(ctx, productID)
}

func (p productService) Update(ctx context.Context, dto dto.UpdateProductDTO) error {
	return p.productStorage.Update(ctx, dto)
}

func (p productService) Approve(ctx context.Context, productID string) error {
	product, err := p.productStorage.GetByID(ctx, productID)
	if err != nil {
		return err
	}
	if product.IsApproved {
		return domain.ErrProductAlreadyApproved
	}
	if err := p.productStorage.Approve(ctx, productID); err != nil {
		return err
	}
	return nil
}

func (p productService) Disapprove(ctx context.Context, productID string) error {
	product, err := p.productStorage.GetByID(ctx, productID)
	if err != nil {
		return err
	}
	if !product.IsApproved {
		return domain.ErrProductAlreadyDisapproved
	}

	if err := p.productStorage.Disapprove(ctx, productID); err != nil {
		return err
	}
	return nil

}
