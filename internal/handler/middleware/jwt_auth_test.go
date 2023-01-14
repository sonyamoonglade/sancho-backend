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

	t.Run("should pass through middleware because token is all valid", func(t *testing.T) {
		app := fiber.New()
		app.Use(m.Use(domain.RoleCustomer))
		app.Get("/ping", func(ctx *fiber.Ctx) error {
			return ctx.Status(http.StatusOK).SendString(pong)
		})

		tokens, err := p.GenerateNewPair(auth.UserAuth{
			UserID: uuid.NewString(),
			Role:   domain.RoleCustomer,
		})
		require.NoError(t, err)

		req := getRequest(tokens.AccessToken)
		res, err := app.Test(req, -1)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, res.StatusCode)
		body := readBody(res.Body)
		require.Equal(t, pong, string(body))
	})

	t.Run("should return 401 Unauthorized because token is missing", func(t *testing.T) {
		app := fiber.New()
		app.Use(m.Use(domain.RoleCustomer))
		app.Get("/ping", func(ctx *fiber.Ctx) error {
			return ctx.Status(http.StatusOK).SendString(pong)
		})
		// Not setting the Authorization header
		req, _ := http.NewRequest(http.MethodGet, baseURL+"/ping", nil)
		res, err := app.Test(req, -1)
		require.NoError(t, err)
		body := readBody(res.Body)
		require.Equal(t, http.StatusUnauthorized, res.StatusCode)
		require.Equal(t, ResponseUnauthorized, string(body))
	})

	t.Run("should return 401 Unauthorized because header is in invalid format", func(t *testing.T) {
		app := fiber.New()
		app.Use(m.Use(domain.RoleCustomer))
		app.Get("/ping", func(ctx *fiber.Ctx) error {
			return ctx.Status(http.StatusOK).SendString(pong)
		})
		// Invalid header format "Bearer ", token is empty
		req := getRequest("")
		res, err := app.Test(req, -1)
		require.NoError(t, err)
		body := readBody(res.Body)
		require.Equal(t, http.StatusUnauthorized, res.StatusCode)
		require.Equal(t, ResponseUnauthorized, string(body))
	})

	t.Run("should return 401 Unauthorized because token is not valid (invalid format)", func(t *testing.T) {
		app := fiber.New()
		app.Use(m.Use(domain.RoleCustomer))
		app.Get("/ping", func(ctx *fiber.Ctx) error {
			return ctx.Status(http.StatusOK).SendString(pong)
		})

		accessToken := "abcderfa"

		req := getRequest(accessToken)
		res, err := app.Test(req, -1)
		require.NoError(t, err)
		body := readBody(res.Body)
		require.Equal(t, http.StatusUnauthorized, res.StatusCode)
		require.Equal(t, ResponseUnauthorized, string(body))
	})

	t.Run("should return 403 Access Denied because role's too low", func(t *testing.T) {
		app := fiber.New()
		// Only for admins!!
		app.Use(m.Use(domain.RoleAdmin))
		app.Get("/ping", func(ctx *fiber.Ctx) error {
			return ctx.Status(http.StatusOK).SendString(pong)
		})

		tokens, err := p.GenerateNewPair(auth.UserAuth{
			Role:   domain.RoleCustomer,
			UserID: uuid.NewString(),
		})

		req := getRequest(tokens.AccessToken)
		res, err := app.Test(req, -1)
		require.NoError(t, err)
		require.Equal(t, http.StatusForbidden, res.StatusCode)
		body := readBody(res.Body)
		require.Equal(t, ResponseAccessDenied, string(body))
	})
}

func getRequest(token string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, baseURL+"/ping", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	return req
}

func readBody(rc io.ReadCloser) []byte {
	body, err := io.ReadAll(rc)
	if err != nil {
		panic(err)
	}
	rc.Close()
	return body
}
