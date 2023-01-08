package dto

import "github.com/sonyamoonglade/sancho-backend/internal/domain"

type CreateProductDTO struct {
	Name         string
	TranslateRU  string
	Description  string
	CategoryName string
	Price        int64
	Features     domain.Features
}

func (d CreateProductDTO) ToDomain() domain.Product {
	return domain.Product{
		Name:        d.Name,
		TranslateRU: d.TranslateRU,
		Description: d.Description,
		ImageURL:    "",
		Price:       d.Price,
		IsApproved:  false,
		Features: domain.Features{
			Weight:      d.Features.Weight,
			Volume:      d.Features.Volume,
			EnergyValue: d.Features.EnergyValue,
			Nutrients:   d.Features.Nutrients,
		},
	}
}

type UpdateProductDTO struct {
	ProductID    string
	Name         *string
	TranslateRU  *string
	Description  *string
	CategoryName *string
	Price        *int64
}
