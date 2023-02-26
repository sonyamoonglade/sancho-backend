package input

import (
	"github.com/sonyamoonglade/sancho-backend/internal/domain"
	"github.com/sonyamoonglade/sancho-backend/internal/services/dto"
)

type CreateUserOrderInput struct {
	Pay             domain.Pay                   `json:"pay" validate:"required"`
	Cart            []CartProductInput           `json:"cart" validate:"required"`
	IsDelivered     bool                         `json:"isDelivered"`
	DeliveryAddress *domain.OrderDeliveryAddress `json:"deliveryAddress,omitempty"`
}

func (c CreateUserOrderInput) ToDTO(customerID string) dto.CreateUserOrderDTO {
	cart := make([]dto.CartProductDTO, 0, len(c.Cart))
	for _, cartProduct := range c.Cart {
		cart = append(cart, dto.CartProductDTO{
			ProductID: cartProduct.ProductID,
			Quantity:  cartProduct.Quantity,
		})
	}
	return dto.CreateUserOrderDTO{
		CustomerID:      customerID,
		Pay:             c.Pay,
		Cart:            cart,
		IsDelivered:     c.IsDelivered,
		DeliveryAddress: c.DeliveryAddress,
	}
}

type CreateWorkerOrderInput struct {
	CustomerName    string                       `json:"customerName" validate:"required"`
	PhoneNumber     string                       `json:"phoneNumber" validate:"required"`
	Cart            []CartProductInput           `json:"cart" validate:"required"`
	DiscountPercent float64                      `json:"discountPercent,omitempty"`
	Pay             domain.Pay                   `json:"pay" validate:"required"`
	DeliveryAddress *domain.OrderDeliveryAddress `json:"deliveryAddress,omitempty"`
	IsDelivered     bool                         `json:"isDelivered"`
}

func (c CreateWorkerOrderInput) ToDTO(customerID string) dto.CreateWorkerOrderDTO {
	cart := make([]dto.CartProductDTO, 0, len(c.Cart))
	for _, cartProduct := range c.Cart {
		cart = append(cart, dto.CartProductDTO{
			ProductID: cartProduct.ProductID,
			Quantity:  cartProduct.Quantity,
		})
	}
	return dto.CreateWorkerOrderDTO{
		CustomerID:      customerID,
		CustomerName:    c.CustomerName,
		PhoneNumber:     c.PhoneNumber,
		Cart:            cart,
		DiscountPercent: c.DiscountPercent,
		Pay:             c.Pay,
		DeliveryAddress: c.DeliveryAddress,
		IsDelivered:     c.IsDelivered,
	}
}

type CartProductInput struct {
	ProductID string `json:"productId" validate:"required"`
	Quantity  int32  `json:"quantity" validate:"required"`
}
