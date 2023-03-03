package service

import (
	"context"
	"errors"
	"math"
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
	orderStorage         storage.Order
	productService       Product
	orderConfig          OrderConfig
	businessMetaProvider domain.MetaProvider
}

func NewOrderService(orderStorage storage.Order,
	productService Product,
	orderConfig OrderConfig,
	metaProvider domain.MetaProvider) Order {
	return &orderService{
		orderStorage:         orderStorage,
		productService:       productService,
		orderConfig:          orderConfig,
		businessMetaProvider: metaProvider,
	}
}

func (o *orderService) GetOrderByID(ctx context.Context, orderID string) (domain.Order, error) {
	return o.orderStorage.GetOrderByID(ctx, orderID)
}

func (o *orderService) GetLastOrderByCustomerID(ctx context.Context, customerID string) (domain.Order, error) {
	return o.orderStorage.GetLastOrderByCustomerID(ctx, customerID)
}

func (o *orderService) GetOrderByNanoIDAt(ctx context.Context, nanoID string, from, to time.Time) (domain.Order, error) {
	return o.orderStorage.GetOrderByNanoIDAt(ctx, nanoID, from, to)
}

// todo: test
func (o *orderService) CreateWorkerOrder(ctx context.Context, dto dto.CreateWorkerOrderDTO) (string, error) {
	amount, cartProducts, err := o.CalculateCartAmount(ctx, dto.Cart)
	if err != nil {
		return "", err
	}

	nanoID, err := o.findNanoID(ctx)
	if err != nil {
		return "", err
	}

	now := time.Now().UTC()
	order := domain.Order{
		NanoID:           nanoID,
		CustomerID:       dto.CustomerID,
		Cart:             cartProducts,
		Pay:              dto.Pay,
		Amount:           amount,
		Discount:         dto.DiscountPercent,
		DiscountedAmount: o.CalculateDiscountedAmount(amount, dto.DiscountPercent),
		Status:           domain.StatusVerified,
		IsDelivered:      dto.IsDelivered,
		DeliveryAddress:  dto.DeliveryAddress,
		CreatedAt:        now,
		VerifiedAt:       &now,
	}

	if dto.IsDelivered && dto.DeliveryAddress != nil {
		meta := o.businessMetaProvider.Get()
		if order.Amount > meta.DeliveryPunishmentThreshold {
			order.DiscountedAmount = o.applyPunishment(order.DiscountedAmount, meta.DeliveryPunishmentValue)
		}
	}

	orderID, err := o.orderStorage.SaveOrder(ctx, order)
	if err != nil {
		return "", err
	}

	return orderID.Hex(), nil
}

func (o *orderService) CreateUserOrder(ctx context.Context, dto dto.CreateUserOrderDTO) (string, error) {
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

	amount, cartProducts, err := o.CalculateCartAmount(ctx, dto.Cart)
	if err != nil {
		return "", err
	}

	nanoID, err := o.findNanoID(ctx)
	if err != nil {
		return "", err
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
		CreatedAt:        now,
	}

	if dto.IsDelivered && dto.DeliveryAddress != nil {
		meta := o.businessMetaProvider.Get()
		if order.Amount > meta.DeliveryPunishmentThreshold {
			order.DiscountedAmount = o.applyPunishment(order.DiscountedAmount, meta.DeliveryPunishmentValue)
		}
	}

	orderID, err := o.orderStorage.SaveOrder(ctx, order)
	if err != nil {
		return "", err
	}

	return orderID.Hex(), nil
}

func (o *orderService) CalculateDiscountedAmount(amount int64, discountPercent float64) int64 {
	return int64(math.Round((1 - discountPercent) * float64(amount)))
}

//todo: test
func (o *orderService) CalculateCartAmount(ctx context.Context, cart []dto.CartProductDTO) (int64, []domain.CartProduct, error) {
	productIDs := make([]string, 0, len(cart))
	for _, product := range cart {
		productIDs = append(productIDs, product.ProductID)
	}

	products, err := o.productService.GetProductsByIDs(ctx, productIDs)
	if err != nil {
		return 0, nil, err
	}

	var (
		total        int64
		cartProducts = make([]domain.CartProduct, 0, len(cart))
	)
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

func (o *orderService) findNanoID(ctx context.Context) (string, error) {
	var (
		day       = time.Hour * 24
		now       = time.Now().UTC()
		yesterday = now.Add(day * -1)
		nanoID    string
		err       error
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

	return nanoID, nil
}

func (o *orderService) applyPunishment(origin, punishment int64) int64 {
	return origin + punishment
}
