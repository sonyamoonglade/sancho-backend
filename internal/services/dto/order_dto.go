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

type CreateWorkerOrderDTO struct {
	CustomerID      string
	CustomerName    string
	PhoneNumber     string
	Cart            []CartProductDTO
	DiscountPercent float64
	Pay             domain.Pay
	DeliveryAddress *domain.OrderDeliveryAddress
	IsDelivered     bool
}

type CartProductDTO struct {
	ProductID string
	Quantity  int32
}
