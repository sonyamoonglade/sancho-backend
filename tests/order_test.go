package tests

import (
	"net/http"
	"testing"
	"time"

	"github.com/sonyamoonglade/sancho-backend/internal/domain"
	"github.com/sonyamoonglade/sancho-backend/internal/handler/input"
)

func (s *APISuite) TestCreateUserOrder() {

	var (
		t       = s.T()
		require = s.Require()
	)

	t.Run("should create user order because customer exists and no limitations applied", func(t *testing.T) {
		cartProduct := products[0].(domain.Product)
		inp := input.CreateUserOrderInput{
			Pay:         domain.PayOnline,
			Cart:        []string{cartProduct.ProductID.Hex()},
			Amount:      cartProduct.Price,
			IsDelivered: true,
			DeliveryAddress: &domain.OrderDeliveryAddress{
				IsAsap:      true,
				Address:     "Орджоникидзе 29а",
				Entrance:    2,
				Floor:       91,
				Apartment:   15,
				DeliveredAt: time.Now().UTC().Add(time.Hour * 1),
			},
		}
		req, _ := http.NewRequest(http.MethodPost, "/api/order/createUserOrder", newBody(inp))
		// customer associated with access token creates an order
		accessToken := newAccessToken(s.tokenProvider, customer.UserID.Hex(), customer.Role)
		req.Header.Set("Authorization", "Bearer "+accessToken)
		res, err := s.app.Test(req, -1)
		require.NoError(err)
		require.Equal(http.StatusCreated, res.StatusCode)
	})

}
