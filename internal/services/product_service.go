package service

import (
	"context"
	"errors"

	"github.com/sonyamoonglade/sancho-backend/internal/appErrors"
	"github.com/sonyamoonglade/sancho-backend/internal/domain"
	"github.com/sonyamoonglade/sancho-backend/internal/services/dto"
	"github.com/sonyamoonglade/sancho-backend/internal/storages"
)

type ProductService struct {
	productStorage storage.Product
}

func NewProductService(productStorage storage.Product) Product {
	return &ProductService{productStorage: productStorage}
}

func (p ProductService) GetByID(ctx context.Context, productID string) (domain.Product, error) {
	product, err := p.productStorage.GetByID(ctx, productID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return domain.Product{}, domain.ErrProductNotFound
		}
		return domain.Product{}, appErrors.WithContext("productStorage.GetByID", err)
	}
	return product, nil
}

func (p ProductService) GetAll(ctx context.Context) ([]domain.Product, error) {
	catalog, err := p.productStorage.GetAll(ctx)
	if err != nil {
		return nil, appErrors.WithContext("productStorage.GetAll", err)
	}
	return catalog, nil
}

func (p ProductService) GetAllCategories(ctx context.Context, sorted bool) ([]domain.Category, error) {
	categories, err := p.productStorage.GetAllCategories(ctx, sorted)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return categories, nil
}

func (p ProductService) Create(ctx context.Context, dto dto.CreateProductDTO) (string, error) {
	category, err := p.productStorage.GetCategoryByName(ctx, dto.CategoryName)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return "", domain.ErrCategoryNotFound
		}
		return "", appErrors.WithContext("productStorage.GetCategoryByName", err)
	}

	product := dto.ToDomain()
	product.Category = category

	productID, err := p.productStorage.Create(ctx, product)
	if err != nil {
		if errors.Is(err, storage.ErrAlreadyExists) {
			return "", domain.ErrProductAlreadyExists
		}

		return "", appErrors.WithContext("productStorage.Create", err)
	}

	return productID.Hex(), nil
}

func (p ProductService) Delete(ctx context.Context, productID string) error {
	if err := p.productStorage.Delete(ctx, productID); err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return domain.ErrProductNotFound
		}
		return appErrors.WithContext("productStorage.Delete", err)
	}
	return nil
}

func (p ProductService) Update(ctx context.Context, dto dto.UpdateProductDTO) error {
	if err := p.productStorage.Update(ctx, dto); err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return domain.ErrProductNotFound
		}
		if updateError, ok := err.(appErrors.UpdateError); ok {
			return updateError
		}
		return appErrors.WithContext("productStorage.Update", err)
	}
	return nil
}

func (p ProductService) Approve(ctx context.Context, productID string) error {
	product, err := p.productStorage.GetByID(ctx, productID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return domain.ErrProductNotFound
		}
		return appErrors.WithContext("productStorage.GetByID", err)
	}
	if product.IsApproved {
		return domain.ErrProductAlreadyApproved
	}
	if err := p.productStorage.Approve(ctx, productID); err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return domain.ErrProductNotFound
		}
		return appErrors.WithContext("productStorage.Approve", err)
	}
	return nil
}

func (p ProductService) Disapprove(ctx context.Context, productID string) error {
	product, err := p.productStorage.GetByID(ctx, productID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return domain.ErrProductNotFound
		}
		return appErrors.WithContext("productStorage.GetByID", err)
	}
	if !product.IsApproved {
		return domain.ErrProductAlreadyDisapproved
	}

	if err := p.productStorage.Disapprove(ctx, productID); err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return domain.ErrProductNotFound
		}
		return appErrors.WithContext("productStorage.Disapprove", err)
	}
	return nil

}
