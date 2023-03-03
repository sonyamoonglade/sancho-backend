package tests

import (
	"context"
	"encoding/json"
	"io"
	"math"
	"net/http"
	"testing"
	"time"

	"github.com/sonyamoonglade/sancho-backend/internal/domain"
	"github.com/sonyamoonglade/sancho-backend/internal/handler/input"
	"github.com/sonyamoonglade/sancho-backend/internal/services/dto"
)

func (s *APISuite) TestCreateUserOrder() {
	var (
		t       = s.T()
		require = s.Require()
	)
	t.Run("should create user order", func(t *testing.T) {
		var (
			cartProduct       = products[0].(domain.Product)
			quantity    int32 = 5
		)
		inp := input.CreateUserOrderInput{
			Pay: domain.PayOnPickup,
			Cart: []input.CartProductInput{
				{ProductID: cartProduct.ProductID.Hex(), Quantity: quantity},
			},
			IsDelivered: false,
		}

		accessToken := newAccessToken(s.tokenProvider, customer.UserID.Hex(), customer.Role)
		req := newRequest("/api/order/create", http.MethodPost, accessToken, newBody(inp))
		res, err := s.app.Test(req, -1)
		printResponseDetails(res)
		require.NoError(err)
		require.Equal(http.StatusCreated, res.StatusCode)

		type createUserOrderResponse struct {
			OrderID string `json:"orderId"`
		}

		var createUserOrderResp createUserOrderResponse
		err = json.NewDecoder(res.Body).Decode(&createUserOrderResp)
		require.NoError(err)

		order, err := s.services.Order.GetOrderByID(context.Background(), createUserOrderResp.OrderID)
		require.NoError(err)

		require.NotZero(len(order.Cart))
		require.True(order.Status == domain.StatusWaitingForVerification)
		require.True(order.CustomerID == customer.UserID.Hex())
		// delivery punishment is not applied because order is not delivered
		require.True(order.DiscountedAmount == order.Amount)
	})
	t.Run("should create user order and apply delivery punishment", func(t *testing.T) {
		var (
			cartProduct       = products[0].(domain.Product)
			quantity    int32 = 5
		)
		inp := input.CreateUserOrderInput{
			Pay: domain.PayOnPickup,
			Cart: []input.CartProductInput{
				{ProductID: cartProduct.ProductID.Hex(), Quantity: quantity},
			},
			IsDelivered: true,
			DeliveryAddress: &domain.OrderDeliveryAddress{
				IsAsap:      true,
				Address:     "blabla",
				Entrance:    1,
				Floor:       2,
				Apartment:   3,
				DeliveredAt: time.Now().UTC().Add(time.Hour * 1),
			},
		}

		accessToken := newAccessToken(s.tokenProvider, customer.UserID.Hex(), customer.Role)
		req := newRequest("/api/order/create", http.MethodPost, accessToken, newBody(inp))
		res, err := s.app.Test(req, -1)
		printResponseDetails(res)
		require.NoError(err)
		require.Equal(http.StatusCreated, res.StatusCode)

		type createUserOrderResponse struct {
			OrderID string `json:"orderId"`
		}

		var createUserOrderResp createUserOrderResponse
		err = json.NewDecoder(res.Body).Decode(&createUserOrderResp)
		require.NoError(err)

		order, err := s.services.Order.GetOrderByID(context.Background(), createUserOrderResp.OrderID)
		require.NoError(err)

		require.NotZero(len(order.Cart))
		require.True(order.Status == domain.StatusWaitingForVerification)
		require.True(order.CustomerID == customer.UserID.Hex())
		// delivery punishment is applied because order is delivered
		require.True(order.DiscountedAmount != order.Amount)
		require.True(order.DiscountedAmount == order.Amount+meta.DeliveryPunishmentValue)
	})

	t.Run("should create user order that is delivered but punishment is not applied", func(t *testing.T) {
		var (
			cartProduct       = products[0].(domain.Product)
			quantity    int32 = 1
		)
		require.True(cartProduct.Price*int64(quantity) <= meta.DeliveryPunishmentThreshold)
		inp := input.CreateUserOrderInput{
			Pay: domain.PayOnPickup,
			Cart: []input.CartProductInput{
				{ProductID: cartProduct.ProductID.Hex(), Quantity: quantity},
			},
			IsDelivered: true,
			DeliveryAddress: &domain.OrderDeliveryAddress{
				IsAsap:      true,
				Address:     "blabla",
				Entrance:    1,
				Floor:       2,
				Apartment:   3,
				DeliveredAt: time.Now().UTC().Add(time.Hour * 1),
			},
		}

		accessToken := newAccessToken(s.tokenProvider, customer.UserID.Hex(), customer.Role)
		req := newRequest("/api/order/create", http.MethodPost, accessToken, newBody(inp))
		res, err := s.app.Test(req, -1)
		printResponseDetails(res)
		require.NoError(err)
		require.Equal(http.StatusCreated, res.StatusCode)

		type createUserOrderResponse struct {
			OrderID string `json:"orderId"`
		}

		var createUserOrderResp createUserOrderResponse
		err = json.NewDecoder(res.Body).Decode(&createUserOrderResp)
		require.NoError(err)

		order, err := s.services.Order.GetOrderByID(context.Background(), createUserOrderResp.OrderID)
		require.NoError(err)

		require.NotZero(len(order.Cart))
		require.True(order.Status == domain.StatusWaitingForVerification)
		require.True(order.CustomerID == customer.UserID.Hex())
		// delivery punishment is not applied because order is delivered but threshold is not reached
		require.True(order.DiscountedAmount == order.Amount)
		require.True(order.Amount <= meta.DeliveryPunishmentThreshold)
	})
}

func (s *APISuite) TestCreateWorkerOrder() {
	var (
		t       = s.T()
		require = s.Require()
	)
	t.Run("should create order for customer that doesn't exist", func(t *testing.T) {
		var (
			cartProduct1       = products[0].(domain.Product)
			cartProduct2       = products[1].(domain.Product)
			quantity     int32 = 5
		)
		inp := input.CreateWorkerOrderInput{
			CustomerName: "Mozart",
			PhoneNumber:  "+79458508374",
			Cart: []input.CartProductInput{
				{ProductID: cartProduct1.ProductID.Hex(), Quantity: quantity},
				{ProductID: cartProduct2.ProductID.Hex(), Quantity: quantity + 2},
			},
			DiscountPercent: 0.0,
			Pay:             domain.PayOnPickup,
			IsDelivered:     false,
			DeliveryAddress: nil,
		}

		accessToken := newAccessToken(s.tokenProvider, worker.UserID.Hex(), worker.Role)
		req := newRequest("/api/order/worker/create", http.MethodPost, accessToken, newBody(inp))
		res, err := s.app.Test(req, -1)
		printResponseDetails(res)
		require.NoError(err)
		require.Equal(http.StatusCreated, res.StatusCode)

		type createWorkerOrderResponse struct {
			OrderID string `json:"orderId"`
		}

		var createWorkerOrderResp createWorkerOrderResponse
		err = json.NewDecoder(res.Body).Decode(&createWorkerOrderResp)
		require.NoError(err)

		order, err := s.services.Order.GetOrderByID(context.Background(), createWorkerOrderResp.OrderID)
		require.NoError(err)

		require.NotZero(len(order.Cart))
		require.True(order.Status == domain.StatusVerified)

		// check if customer was registered and order is created for him
		registeredCustomer, err := s.services.User.GetCustomerByPhoneNumber(context.Background(), inp.PhoneNumber)
		require.NoError(err)
		require.Equal(inp.PhoneNumber, registeredCustomer.PhoneNumber)
		require.Equal(inp.CustomerName, *registeredCustomer.Name)
		require.Equal(order.CustomerID, registeredCustomer.UserID.Hex())
	})

	t.Run("should create order with discount for customer that exists", func(t *testing.T) {
		var (
			cartProduct1       = products[0].(domain.Product)
			cartProduct2       = products[1].(domain.Product)
			quantity     int32 = 5
		)
		inp := input.CreateWorkerOrderInput{
			CustomerName: *customer.Name,
			PhoneNumber:  customer.PhoneNumber,
			Cart: []input.CartProductInput{
				{ProductID: cartProduct1.ProductID.Hex(), Quantity: quantity},
				{ProductID: cartProduct2.ProductID.Hex(), Quantity: quantity + 2},
			},
			DiscountPercent: 0.15,
			Pay:             domain.PayOnPickup,
			IsDelivered:     false,
			DeliveryAddress: nil,
		}

		accessToken := newAccessToken(s.tokenProvider, worker.UserID.Hex(), worker.Role)
		req := newRequest("/api/order/worker/create", http.MethodPost, accessToken, newBody(inp))
		res, err := s.app.Test(req, -1)
		printResponseDetails(res)
		require.NoError(err)
		require.Equal(http.StatusCreated, res.StatusCode)

		type createWorkerOrderResponse struct {
			OrderID string `json:"orderId"`
		}

		var createWorkerOrderResp createWorkerOrderResponse
		err = json.NewDecoder(res.Body).Decode(&createWorkerOrderResp)
		require.NoError(err)

		order, err := s.services.Order.GetOrderByID(context.Background(), createWorkerOrderResp.OrderID)
		require.NoError(err)

		require.NotZero(len(order.Cart))
		require.True(order.Status == domain.StatusVerified)
		require.Equal(order.CustomerID, customer.UserID.Hex())
		validAmount := int64(math.Round(float64(order.Amount) * (1 - order.Discount)))
		require.Equal(validAmount, order.DiscountedAmount)
	})

	t.Run("should create worker order and apply delivery punishment that was set", func(t *testing.T) {
		var (
			cartProduct1 = products[0].(domain.Product)
			cartProduct2 = products[1].(domain.Product)
			quantity     = int32(5)
		)
		inp := input.CreateWorkerOrderInput{
			CustomerName: *customer.Name,
			PhoneNumber:  customer.PhoneNumber,
			Cart: []input.CartProductInput{
				{ProductID: cartProduct1.ProductID.Hex(), Quantity: quantity},
				{ProductID: cartProduct2.ProductID.Hex(), Quantity: quantity + 2},
			},
			DiscountPercent: 0.15,
			Pay:             domain.PayOnPickup,
			IsDelivered:     true,
			DeliveryAddress: &domain.OrderDeliveryAddress{
				IsAsap:      true,
				Address:     "random-address",
				Entrance:    2,
				Floor:       3,
				Apartment:   4,
				DeliveredAt: time.Now().UTC().Add(time.Hour * 1),
			},
		}

		accessToken := newAccessToken(s.tokenProvider, worker.UserID.Hex(), worker.Role)
		req := newRequest("/api/order/worker/create", http.MethodPost, accessToken, newBody(inp))
		res, err := s.app.Test(req, -1)
		printResponseDetails(res)
		require.NoError(err)
		require.Equal(http.StatusCreated, res.StatusCode)

		type createWorkerOrderResponse struct {
			OrderID string `json:"orderId"`
		}

		var createWorkerOrderResp createWorkerOrderResponse
		err = json.NewDecoder(res.Body).Decode(&createWorkerOrderResp)
		require.NoError(err)

		order, err := s.services.Order.GetOrderByID(context.Background(), createWorkerOrderResp.OrderID)
		require.NoError(err)

		// prepare cart
		var cart []dto.CartProductDTO
		for _, cartItem := range inp.Cart {
			cart = append(cart, dto.CartProductDTO{
				ProductID: cartItem.ProductID,
				Quantity:  cartItem.Quantity,
			})
		}

		cartAmount, _, err := s.services.Order.CalculateCartAmount(context.Background(), cart)
		require.NoError(err)

		discountedAmount := s.services.Order.CalculateDiscountedAmount(cartAmount, inp.DiscountPercent)

		// because delivery punishment is applied (amount is more than meta.DeliveryPunishmentThreshold)
		require.NotEqual(order.DiscountedAmount, discountedAmount)
		require.True(order.DiscountedAmount == discountedAmount+meta.DeliveryPunishmentValue)
	})
}

func newRequest(url, method, token string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	return req
}
