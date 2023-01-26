package service

import (
	"context"
	"errors"
	"time"

	"github.com/sonyamoonglade/sancho-backend/internal/domain"
	"github.com/sonyamoonglade/sancho-backend/internal/services/dto"
	storage "github.com/sonyamoonglade/sancho-backend/internal/storages"
	"github.com/sonyamoonglade/sancho-backend/pkg/nanoid"
)

type OrderConfig struct {
	// Duration in minutes that represents minimal time to wait in order to create
	// new Order when having pending Order (status=waiting for verification)
	PendingOrderWaitTime time.Duration
}

type orderService struct {
	orderStorage   storage.Order
	productService Product
	orderConfig    OrderConfig
}

func NewOrderService(orderStorage storage.Order, productService Product, orderConfig OrderConfig) Order {
	return &orderService{orderStorage: orderStorage, productService: productService, orderConfig: orderConfig}
}

func (o orderService) GetOrderByID(ctx context.Context, orderID string) (domain.Order, error) {
	return o.orderStorage.GetOrderByID(ctx, orderID)
}

func (o orderService) GetLastOrderByCustomerID(ctx context.Context, customerID string) (domain.Order, error) {
	return o.orderStorage.GetLastOrderByCustomerID(ctx, customerID)
}

func (o orderService) GetOrderByNanoIDAt(ctx context.Context, nanoID string, from, to time.Time) (domain.Order, error) {
	return o.orderStorage.GetOrderByNanoIDAt(ctx, nanoID, from, to)
}

// todo: test
func (o orderService) CreateUserOrder(ctx context.Context, dto dto.CreateUserOrderDTO) (string, error) {
	// Firstly, check for pending order
	pendingOrder, err := o.GetLastOrderByCustomerID(ctx, dto.CustomerID)
	if err != nil && !errors.Is(err, domain.ErrOrderNotFound) {
		return "", err
	}
	var (
		now                  = time.Now().UTC()
		pendingOrderWaitTime = o.orderConfig.PendingOrderWaitTime
		canCreateNewOrder    = pendingOrder.CreatedAt.Add(pendingOrderWaitTime).Before(now)
	)
	// Do not allow users to create another order
	// when one's pending (waiting for verification) and wait time has not passed yet
	if pendingOrder.Status == domain.StatusWaitingForVerification && !canCreateNewOrder {
		return "", domain.ErrHavePendingOrder
	}

	amount, cartProducts, err := o.calculateCartAmount(ctx, dto.Cart)
	if err != nil {
		return "", err
	}

	var (
		nanoID    string
		day       = time.Hour * 24
		yesterday = now.Add(day * -1)
	)
	for {
		nanoID, err = nanoid.GenerateNanoID()
		if err != nil {
			return "", err
		}

		// Look for orders within 24h to have same nanoID. It's looking in [now -24h, now]
		_, err := o.GetOrderByNanoIDAt(ctx, nanoID, yesterday, now)
		if err != nil {
			// Non-duplicate nanoID has found so we stop
			if errors.Is(err, domain.ErrOrderNotFound) {
				break
			}
			// Internal error
			return "", err
		}
	}

	order := domain.Order{
		NanoID:     nanoID,
		CustomerID: dto.CustomerID,
		Cart:       cartProducts,
		Pay:        dto.Pay,
		Amount:     amount,
		// If customer creates an order it can't get any discount
		Discount:         0,
		DiscountedAmount: amount,
		Status:           domain.StatusWaitingForVerification,
		IsDelivered:      dto.IsDelivered,
		DeliveryAddress:  dto.DeliveryAddress,
		CreatedAt:        time.Now().UTC(),
	}

	orderID, err := o.orderStorage.SaveOrder(ctx, order)
	if err != nil {
		return "", err
	}

	return orderID.Hex(), nil
}

//todo: test
func (o orderService) calculateCartAmount(ctx context.Context, cart []dto.CartProductDTO) (int64, []domain.CartProduct, error) {
	productIDs := make([]string, 0, len(cart))
	for _, product := range cart {
		productIDs = append(productIDs, product.ProductID)
	}

	products, err := o.productService.GetProductsByIDs(ctx, productIDs)
	if err != nil {
		return 0, nil, err
	}

	var total int64
	cartProducts := make([]domain.CartProduct, 0, len(cart))
	for _, cartProduct := range cart {
		for _, product := range products {
			if cartProduct.ProductID == product.ProductID.Hex() {
				total += product.Price * int64(cartProduct.Quantity)
				cartProducts = append(cartProducts, domain.CartProduct{
					Product:  product,
					Quantity: cartProduct.Quantity,
				})
			}
		}
	}
	return total, cartProducts, nil
}
