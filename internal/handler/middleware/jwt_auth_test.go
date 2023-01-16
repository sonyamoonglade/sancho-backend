package middleware

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/sonyamoonglade/sancho-backend/internal/domain"
	mock_service "github.com/sonyamoonglade/sancho-backend/internal/services/mocks"
	"github.com/sonyamoonglade/sancho-backend/pkg/auth"
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

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	p, err := auth.NewProvider(ttl, key, issuer)
	require.NoError(t, err)

	authService := mock_service.NewMockAuth(ctrl)
	m := NewJWTAuthMiddleware(authService, p)
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

	t.Run("should return 401 because access token has expired and no refresh token is present in cookies", func(t *testing.T) {
		app := fiber.New()
		app.Use(m.Use(domain.RoleCustomer))
		app.Get("/ping", func(ctx *fiber.Ctx) error {
			return ctx.Status(http.StatusOK).SendString(pong)
		})

		tokens, err := p.GenerateNewPairWithTTL(auth.UserAuth{
			UserID: uuid.NewString(),
			Role:   domain.RoleCustomer,
		}, time.Millisecond*5)
		require.NoError(t, err)
		time.Sleep(time.Millisecond * 5)

		req := getRequest(tokens.AccessToken)
		res, err := app.Test(req, -1)
		require.NoError(t, err)
		require.Equal(t, http.StatusUnauthorized, res.StatusCode)
		body := readBody(res.Body)
		require.Equal(t, ResponseTokenExpired, string(body))
	})

	t.Run("should return 200 OK and rotate tokens because access token has expired but valid refresh token is present", func(t *testing.T) {
		var userID = uuid.NewString()
		app := fiber.New()
		app.Use(m.Use(domain.RoleAdmin))
		app.Get("/ping", func(ctx *fiber.Ctx) error {
			return ctx.Status(http.StatusOK).SendString(pong)
		})

		tokens, err := p.GenerateNewPairWithTTL(auth.UserAuth{
			UserID: userID,
			Role:   domain.RoleAdmin,
		}, time.Millisecond*5)
		require.NoError(t, err)
		time.Sleep(time.Millisecond * 5)

		req := getRequest(tokens.AccessToken)
		req.AddCookie(&http.Cookie{
			Name:     "refresh-token",
			Value:    tokens.RefreshToken,
			MaxAge:   60 * 60 * 24 * 30,
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		})

		// Valid token pair. It will be returned by mock service.
		newTokens, _ := p.GenerateNewPairWithTTL(auth.UserAuth{
			UserID: userID,
			Role:   domain.RoleAdmin,
		}, time.Millisecond*10)

		authService.EXPECT().
			RefreshAdminToken(gomock.Any(), userID, tokens.RefreshToken).
			Return(newTokens, nil).
			Times(1)

		res, err := app.Test(req, -1)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, res.StatusCode)

		type response struct {
			AccessToken string `json:"accessToken"`
		}
		body := readBody(res.Body)
		var r response
		err = json.Unmarshal(body, &r)
		require.NoError(t, err)

		require.Equal(t, newTokens.AccessToken, r.AccessToken)
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
