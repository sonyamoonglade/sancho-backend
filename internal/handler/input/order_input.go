package input

import (
	"github.com/sonyamoonglade/sancho-backend/internal/domain"
	"github.com/sonyamoonglade/sancho-backend/internal/services/dto"
)

type CreateUserOrderInput struct {
	Pay             domain.Pay                   `json:"pay" validate:"required"`
	Cart            []CartProductInput           `json:"cart" validate:"required"`
	IsDelivered     bool                         `json:"isDelivered" validate:"required"`
	DeliveryAddress *domain.OrderDeliveryAddress `json:"deliveryAddress,omitempty"`
}

type CartProductInput struct {
	ProductID string `json:"productId" validate:"required"`
	Quantity  int32  `json:"quantity" validate:"required"`
}

func (cu CreateUserOrderInput) ToDTO(customerID string) dto.CreateUserOrderDTO {
	cart := make([]dto.CartProductDTO, 0, len(cu.Cart))
	for _, cartProduct := range cu.Cart {
		cart = append(cart, dto.CartProductDTO{
			ProductID: cartProduct.ProductID,
			Quantity:  cartProduct.Quantity,
		})
	}
	return dto.CreateUserOrderDTO{
		CustomerID:      customerID,
		Pay:             cu.Pay,
		Cart:            cart,
		IsDelivered:     cu.IsDelivered,
		DeliveryAddress: cu.DeliveryAddress,
	}
}
