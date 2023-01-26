package tests

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/sonyamoonglade/sancho-backend/internal/domain"
	"github.com/sonyamoonglade/sancho-backend/internal/handler/input"
)

func (s *APISuite) TestCreateUserOrder() {

	var (
		t       = s.T()
		require = s.Require()
	)

	t.Run("should create user order because customer exists and no limitations applied", func(t *testing.T) {
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
		req := newRequest("/api/order/createUserOrder", http.MethodPost, accessToken, newBody(inp))
		res, err := s.app.Test(req, -1)
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
