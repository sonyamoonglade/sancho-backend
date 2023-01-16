package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/sonyamoonglade/sancho-backend/internal/domain"
	service "github.com/sonyamoonglade/sancho-backend/internal/services"
	"github.com/sonyamoonglade/sancho-backend/pkg/auth"
)

const (
	HeaderAuthorization = "Authorization"

	ResponseUnauthorized = "unauthorized"
	ResponseTokenExpired = "token has expired"
	ResponseAccessDenied = "access denied"
)

// JWTAuthMiddleware validates incoming Bearer access token,
// does RBAC, checking incoming user role against required
// and enriches request's context with userID of the token owner.
type JWTAuthMiddleware struct {
	tokenProvider auth.TokenProvider
	authService   service.Auth
}

func NewJWTAuthMiddleware(authService service.Auth, tokenProvider auth.TokenProvider) *JWTAuthMiddleware {
	return &JWTAuthMiddleware{
		authService:   authService,
		tokenProvider: tokenProvider,
	}
}

func (m JWTAuthMiddleware) Use(requiredRole domain.Role) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var (
			headers    = c.GetReqHeaders()
			authHeader = headers[HeaderAuthorization]
			token      string
		)
		if authHeader == "" {
			return c.Status(http.StatusUnauthorized).SendString(ResponseUnauthorized)
		}
		split := strings.Split(authHeader, " ")
		if len(split) == 1 || split[1] == "" {
			return c.Status(http.StatusUnauthorized).SendString(ResponseUnauthorized)
		}

		token = split[1]
		userAuth, err := m.tokenProvider.ParseAndValidate(token)
		// Token isn't valid
		if err != nil {
			switch true {
			// If token has expired userAuth is not nil. See jwt.go:122
			case errors.Is(err, auth.ErrTokenExpired):
				return m.refreshTokens(c, userAuth)
			case errors.Is(err, auth.ErrInvalidToken),
				errors.Is(err, auth.ErrInvalidIssuer):
				return c.Status(http.StatusUnauthorized).SendString(ResponseUnauthorized)
			}
			// Delegate internal errors to ErrorHandler
			return err
		}

		// Check if role's insufficient.
		if !userAuth.Role.CheckPermissions(requiredRole) {
			return c.Status(http.StatusForbidden).SendString(ResponseAccessDenied)
		}

		c.Locals(userAuth, userAuth.UserID)
		return c.Next()
	}
}

func (m JWTAuthMiddleware) refreshTokens(c *fiber.Ctx, userAuth auth.UserAuth) error {
	// Refresh tokens here
	refreshToken := c.Cookies("refresh-token", "")
	if refreshToken == "" {
		return c.Status(http.StatusUnauthorized).SendString(ResponseTokenExpired)
	}

	switch userAuth.Role {
	case domain.RoleCustomer:
		//TODO
		return nil
	case domain.RoleWorker:
		// TODO
		return nil
	case domain.RoleAdmin:
		newTokens, err := m.authService.RefreshAdminToken(c.Context(), userAuth.UserID, refreshToken)
		if err != nil {
			return err
		}
		m.setRefreshTokenCookie(c, newTokens.RefreshToken)
		return m.returnAccessToken(c, newTokens.AccessToken)
	default:
		return c.SendStatus(http.StatusInternalServerError)
	}
}

func (m JWTAuthMiddleware) setRefreshTokenCookie(c *fiber.Ctx, refreshToken string) {
	c.Cookie(&fiber.Cookie{
		Name:  "refresh-token",
		Value: refreshToken,
		// 30 Days
		MaxAge:   60 * 60 * 24 * 30,
		Secure:   true,
		HTTPOnly: true,
		SameSite: "lax",
	})
}

func (m JWTAuthMiddleware) returnAccessToken(c *fiber.Ctx, accessToken string) error {
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"accessToken": accessToken,
	})
}
