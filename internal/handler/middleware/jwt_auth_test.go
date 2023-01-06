package middleware

import (
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sonyamoonglade/sancho-backend/auth"
	"github.com/sonyamoonglade/sancho-backend/internal/domain"
	"github.com/stretchr/testify/require"
)

const baseURL = "http://localhost"

func TestJWTAuthMiddleware(t *testing.T) {
	var (
		ttl    = time.Millisecond * 100
		key    = []byte("abcd")
		issuer = baseURL
		pong   = "pong"
	)

	p, err := auth.NewProvider(ttl, key, issuer)
	require.NoError(t, err)
	m := NewJWTAuthMiddleware(p)
	require.NotNil(t, m)

	app := fiber.New()
	t.Run("should pass into customer endpoint and return pong", func(t *testing.T) {
		app.Use(m.Use(domain.RoleCustomer))
		app.Get("/ping", func(ctx *fiber.Ctx) error {
			return ctx.Status(http.StatusOK).SendString(pong)
		})

		tokens, err := p.GenerateNewPair(auth.UserAuth{
			UserID: uuid.NewString(),
			Role:   domain.RoleCustomer,
		})
		require.NoError(t, err)

		req, _ := http.NewRequest(http.MethodGet, baseURL+"/ping", nil)
		req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)
		res, err := app.Test(req)
		require.NoError(t, err)

		require.Equal(t, http.StatusOK, res.StatusCode)
		body, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		defer res.Body.Close()
		require.Equal(t, pong, string(body))
	})
}
