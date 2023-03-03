package service

import (
	"context"
	"math"
	"testing"
	"time"

	f "github.com/brianvoe/gofakeit/v6"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/sonyamoonglade/sancho-backend/internal/domain"
	"github.com/sonyamoonglade/sancho-backend/internal/services/dto"
	mock_service "github.com/sonyamoonglade/sancho-backend/internal/services/mocks"
	mock_storage "github.com/sonyamoonglade/sancho-backend/internal/storages/mocks"
	"github.com/sonyamoonglade/sancho-backend/pkg/meta_cache"
	"github.com/sonyamoonglade/sancho-backend/pkg/nanoid"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCreateOrder(t *testing.T) {
	t.Run("should not create order because customer has pending order", func(t *testing.T) {
		orderService, productService, orderStorage := getServices(t, OrderConfig{
			PendingOrderWaitTime: time.Minute * 5,
		})
		_ = productService
		mockProduct := getProduct()
		quantity := int32(5)
		d := dto.CreateUserOrderDTO{
			CustomerID: primitive.NewObjectID().Hex(),
			Pay:        domain.PayOnPickup,
			Cart: []dto.CartProductDTO{
				{ProductID: mockProduct.ProductID.Hex(), Quantity: quantity},
			},
			IsDelivered:     false,
			DeliveryAddress: nil,
		}
		var (
			// In order to return ErrHavePending order createdAt of pending order + PendingOrderWaitTime
			// should be greater than current time. Assuming that PendingOrderWaitTime is 5 minutes, see line:20
			// createdAt should be between (now-5minutes,now)
			mockCreatedAt = time.Now().UTC().Add(time.Minute * -2)
			mockCart      = []domain.CartProduct{{
				Product:  mockProduct,
				Quantity: 2,
			}}
			// Pending orders that have status waiting for verification
			mockStatus       = domain.StatusWaitingForVerification
			mockPendingOrder = getOrder(d.CustomerID, mockCreatedAt, mockCart, mockStatus)
		)

		orderStorage.EXPECT().GetLastOrderByCustomerID(gomock.Any(), d.CustomerID).Return(mockPendingOrder, nil)

		orderID, err := orderService.CreateUserOrder(context.Background(), d)
		require.Error(t, err)
		require.Zero(t, orderID)
		require.Equal(t, domain.ErrHavePendingOrder, err)
	})
	t.Run("should create order because customer has no pending order", func(t *testing.T) {
		orderService, productService, orderStorage := getServices(t, OrderConfig{
			PendingOrderWaitTime: time.Minute * 5,
		})
		mockProduct := getProduct()
		quantity := int32(5)
		d := dto.CreateUserOrderDTO{
			CustomerID: primitive.NewObjectID().Hex(),
			Pay:        domain.PayOnPickup,
			Cart: []dto.CartProductDTO{
				{ProductID: mockProduct.ProductID.Hex(), Quantity: quantity},
			},
			IsDelivered:     false,
			DeliveryAddress: nil,
		}
		var (
			mockCreatedAt = time.Now().UTC().Add(time.Minute * -6)
			mockCart      = []domain.CartProduct{{
				Product:  mockProduct,
				Quantity: 2,
			}}
			// Pending orders that have status waiting for verification
			mockStatus = domain.StatusWaitingForVerification
			// The order is not technically pending, because mockCreatedAt is over PendingOrderWaitTime, so
			// it's fine to create another order aside the mockPendingOrder.
			mockPendingOrder = getOrder(d.CustomerID, mockCreatedAt, mockCart, mockStatus)
		)

		orderStorage.
			EXPECT().
			GetLastOrderByCustomerID(gomock.Any(), d.CustomerID).
			Return(mockPendingOrder, nil)

		productService.
			EXPECT().
			GetProductsByIDs(gomock.Any(), []string{d.Cart[0].ProductID}).
			Return([]domain.Product{mockProduct}, nil)

		// Assume that no order has been created within 24h with such nanoID generated inside CreateUserOrder method
		orderStorage.
			EXPECT().
			GetOrderByNanoIDAt(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(domain.Order{}, domain.ErrOrderNotFound)

		orderStorage.
			EXPECT().
			SaveOrder(gomock.Any(), gomock.AssignableToTypeOf(domain.Order{})).
			DoAndReturn(func(ctx context.Context, order domain.Order) (primitive.ObjectID, error) {
				// Make some assertions here in order to assume that order saved is correct
				require.Equal(t, d.CustomerID, order.CustomerID)
				require.Equal(t, d.Pay, order.Pay)
				require.EqualValues(t, d.DeliveryAddress, order.DeliveryAddress)
				require.Equal(t, d.IsDelivered, order.IsDelivered)
				// Users cant get discount on an order when created by themselves
				require.True(t, order.Discount == 0)
				require.True(t, order.Amount == order.DiscountedAmount)
				return primitive.NewObjectID(), nil
			}).
			Times(1)

		orderID, err := orderService.CreateUserOrder(context.Background(), d)
		require.NoError(t, err)
		require.NotZero(t, orderID)
	})
}

func TestCalculateCartAmount(t *testing.T) {

	t.Run("should calculate amount of products returned by productService", func(t *testing.T) {
		orderService, productService, orderStorage := getServices(t, OrderConfig{
			PendingOrderWaitTime: time.Minute * 5,
		})
		_ = orderStorage

		var (
			products   []domain.Product
			productIDs []string
			cart       []dto.CartProductDTO
			total      int64
		)
		for i := 0; i < 5; i++ {
			p := getProduct()
			products = append(products, p)
			productIDs = append(productIDs, p.ProductID.Hex())

			quantity := int32(f.IntRange(1, 10))
			cart = append(cart, dto.CartProductDTO{
				ProductID: p.ProductID.Hex(),
				Quantity:  quantity,
			})

			// Calculate total price of all products in cart
			total += p.Price * int64(quantity)
		}

		productService.
			EXPECT().
			GetProductsByIDs(gomock.Any(), productIDs).
			Return(products, nil).
			Times(1)

		amount, cartProducts, err := orderService.CalculateCartAmount(context.Background(), cart)
		require.NoError(t, err)
		require.NotNil(t, cartProducts)
		require.Equal(t, total, amount)

		// Check if all products were added to cartProducts with the same quantity and original product values
		for _, returnedCartProduct := range cartProducts {
			for _, cartProduct := range cart {
				for _, product := range products {
					var (
						sameProductInCart                 = product.ProductID.Hex() == cartProduct.ProductID
						sameProductInReturnedCartProducts = product.ProductID.Hex() == returnedCartProduct.ProductID.Hex()
					)
					if sameProductInCart && sameProductInReturnedCartProducts {
						// Quantity is correctly copied
						require.Equal(t, cartProduct.Quantity, returnedCartProduct.Quantity)
						// Product that's added is correctly copied
						require.EqualValues(t, returnedCartProduct.Product, product)
					}
				}
			}
		}
	})
	t.Run("should return 0 because no such products exist", func(t *testing.T) {
		orderService, productService, orderStorage := getServices(t, OrderConfig{
			PendingOrderWaitTime: time.Minute * 5,
		})
		_ = orderStorage

		var (
			products   []domain.Product
			productIDs []string
			cart       []dto.CartProductDTO
		)
		for i := 0; i < 5; i++ {
			p := getProduct()
			products = append(products, p)
			fakeCartID := uuid.NewString()
			productIDs = append(productIDs, fakeCartID)

			quantity := int32(f.IntRange(1, 10))

			// Append random productID's to cart in order to fake that products do not exist
			cart = append(cart, dto.CartProductDTO{
				ProductID: fakeCartID,
				Quantity:  quantity,
			})

		}

		productService.
			EXPECT().
			GetProductsByIDs(gomock.Any(), productIDs).
			Return(nil, domain.ErrNoProducts).
			Times(1)

		amount, cartProducts, err := orderService.CalculateCartAmount(context.Background(), cart)
		require.Error(t, err)
		require.Equal(t, domain.ErrNoProducts, err)
		require.Nil(t, cartProducts)
		require.Zero(t, amount)
	})
}

func TestCalculateDiscountedAmount(t *testing.T) {

	orderService, productService, orderStorage := getServices(t, OrderConfig{
		PendingOrderWaitTime: time.Minute * 5,
	})
	_ = orderStorage
	_ = productService

	tests := []struct {
		amount          int64
		discountPercent float64
		expected        int64
	}{
		{1000, 0.5, 500},
		{1000, 0.6, 400},
		{1234, 0.15, 1049},
		{381, 0.05, 362},
	}

	for _, test := range tests {
		require.Equal(t, test.expected, orderService.CalculateDiscountedAmount(test.amount, test.discountPercent))
	}

}

func getServices(t *testing.T, orderConfig OrderConfig) (*orderService, *mock_service.MockProduct, *mock_storage.MockOrder) {
	ctrl := gomock.NewController(t)
	orderStorage := mock_storage.NewMockOrder(ctrl)
	productService := mock_service.NewMockProduct(ctrl)
	metaCache := meta_cache.NewMetaCache()
	metaCache.Set(domain.BusinessMeta{
		DeliveryPunishmentThreshold: 400,
		DeliveryPunishmentValue:     100,
	})
	ordService := NewOrderService(orderStorage, productService, orderConfig, metaCache)
	return ordService.(*orderService), productService, orderStorage
}

func getOrder(customerID string, createdAt time.Time, cart []domain.CartProduct, orderStatus domain.OrderStatus) domain.Order {
	nanoId, err := nanoid.GenerateNanoID()
	if err != nil {
		panic(err)
	}
	var (
		amount   = f.Int64()
		discount = f.Float64Range(0.0, 1.0)
	)
	return domain.Order{
		OrderID:           primitive.NewObjectID(),
		NanoID:            nanoId,
		CustomerID:        customerID,
		Cart:              cart,
		Pay:               domain.PayOnPickup,
		Amount:            amount,
		Discount:          discount,
		DiscountedAmount:  int64(math.Round((1 - discount) * float64(amount))),
		Status:            orderStatus,
		IsDelivered:       false,
		DeliveryAddress:   nil,
		CreatedAt:         createdAt,
		CancelExplanation: nil,
	}
}

func getProduct() domain.Product {
	return domain.Product{
		ProductID:   primitive.NewObjectID(),
		Name:        f.BeerName(),
		TranslateRU: f.HipsterWord(),
		Description: f.LoremIpsumSentence(10),
		ImageURL:    stringPtr(f.ImageURL(200, 300)),
		IsApproved:  f.Bool(),
		Price:       f.Int64(),
		Category: domain.Category{
			CategoryID: primitive.NewObjectID(),
			Rank:       int32(f.IntRange(1, 10)),
			Name:       f.HipsterWord(),
		},
		Features: domain.Features{
			IsLiquid:    f.Bool(),
			Weight:      f.Int32(),
			Volume:      f.Int32(),
			EnergyValue: f.Int32(),
			Nutrients: &domain.Nutrients{
				Carbs:    f.Int32(),
				Proteins: f.Int32(),
				Fats:     f.Int32(),
			},
		},
	}
}

func stringPtr(s string) *string {
	return &s
}
