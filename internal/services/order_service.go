package service

import (
	"context"
	"time"

	"github.com/sonyamoonglade/sancho-backend/internal/domain"
	"github.com/sonyamoonglade/sancho-backend/internal/services/dto"
	storage "github.com/sonyamoonglade/sancho-backend/internal/storages"
	"github.com/sonyamoonglade/sancho-backend/pkg/nanoid"
)

type orderService struct {
	orderStorage storage.Order
}

func NewOrderService(orderStorage storage.Order) Order {
	return &orderService{orderStorage: orderStorage}
}
func (o orderService) CreateUserOrder(ctx context.Context, dto dto.CreateUserOrderDTO) (string, error) {
	nanoID, err := nanoid.GenerateNanoID()
	if err != nil {
		return "", err
	}
	//todo:
	//dto.Cart

	order := domain.Order{
		NanoID:     nanoID,
		CustomerID: dto.CustomerID,
		Cart:       nil,
		Pay:        dto.Pay,
		Amount:     0,
		// If customer creates an order it can't get any discount
		Discount:         0,
		DiscountedAmount: 0,
		Status:           domain.StatusWaitingForVerification,
		IsDelivered:      dto.IsDelivered,
		DeliveryAddress:  dto.DeliveryAddress,
		CreatedAt:        time.Now().UTC(),
	}
	_ = order
	return "", nil

}
