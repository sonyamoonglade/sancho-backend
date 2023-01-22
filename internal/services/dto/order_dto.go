package dto

import (
	"github.com/sonyamoonglade/sancho-backend/internal/domain"
)

type CreateUserOrderDTO struct {
	CustomerID      string
	Pay             domain.Pay
	Cart            []CartProductDTO
	IsDelivered     bool
	DeliveryAddress *domain.OrderDeliveryAddress
}

type CartProductDTO struct {
	ProductID string
	Quantity  int32
}
