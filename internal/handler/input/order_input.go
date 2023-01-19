package input

import (
	"github.com/sonyamoonglade/sancho-backend/internal/domain"
	"github.com/sonyamoonglade/sancho-backend/internal/services/dto"
)

type CreateUserOrderInput struct {
	Pay             domain.Pay                   `json:"pay"`
	Cart            []string                     `json:"cart"`
	Amount          int64                        `json:"amount"`
	IsDelivered     bool                         `json:"isDelivered"`
	DeliveryAddress *domain.OrderDeliveryAddress `json:"deliveryAddress,omitempty"`
}

func (cu CreateUserOrderInput) ToDTO(customerID string) dto.CreateUserOrderDTO {
	return dto.CreateUserOrderDTO{
		CustomerID:      customerID,
		Pay:             cu.Pay,
		Cart:            cu.Cart,
		Amount:          cu.Amount,
		IsDelivered:     cu.IsDelivered,
		DeliveryAddress: cu.DeliveryAddress,
	}
}
