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

func (p ProductService) Create(ctx context.Context, dto dto.CreateProductDTO) error {
	category, err := p.productStorage.GetCategoryByName(ctx, dto.CategoryName)
	if err != nil {
		return appErrors.WithContext("productStorage.GetCategoryByName", err)
	}

	product := dto.ToDomain()
	product.Category = category

	if err := p.productStorage.Create(ctx, product); err != nil {
		if errors.Is(err, storage.ErrAlreadyExists) {
			return domain.ErrProductAlreadyExists
		}

		return appErrors.WithContext("productStorage.Create", err)
	}

	return nil
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

func (p ProductService) Update(ctx context.Context, productID string, dto dto.UpdateProductDTO) error {
	//TODO implement me
	panic("implement me")
}

func (p ProductService) Approve(ctx context.Context, productID string) error {
	if err := p.productStorage.Approve(ctx, productID); err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return domain.ErrProductNotFound
		}
		return appErrors.WithContext("productStorage.Approve", err)
	}
	return nil
}

func (p ProductService) ChangeImageURL(ctx context.Context, productID string, imageURL string) error {
	if err := p.productStorage.ChangeImageURL(ctx, productID, imageURL); err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return domain.ErrProductNotFound
		}
		return appErrors.WithContext("productStorage.ChangeImageURL", err)
	}
	return nil
}
