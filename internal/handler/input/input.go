package input

import (
	"github.com/sonyamoonglade/sancho-backend/internal/domain"
	"github.com/sonyamoonglade/sancho-backend/internal/services/dto"
)

type CreateProductInput struct {
	Name         string          `json:"name" validate:"required"`
	TranslateRU  string          `json:"translateRU" validate:"required"`
	Description  string          `json:"description" validate:"required"`
	CategoryName string          `json:"categoryName" validate:"required"`
	Price        int64           `json:"price" validate:"required"`
	Features     domain.Features `json:"features" validate:"required"`
}

func (c CreateProductInput) ToDTO() dto.CreateProductDTO {
	return dto.CreateProductDTO{
		Name:         c.Name,
		TranslateRU:  c.TranslateRU,
		Description:  c.Description,
		CategoryName: c.CategoryName,
		Price:        c.Price,
		Features:     c.Features,
	}
}
