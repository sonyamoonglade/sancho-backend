package dto

import "github.com/sonyamoonglade/sancho-backend/internal/domain"

type CreateUserOrderDTO struct {
	CustomerID      string
	Pay             domain.Pay
	Cart            []string
	Amount          int64
	IsDelivered     bool
	DeliveryAddress *domain.OrderDeliveryAddress
}
