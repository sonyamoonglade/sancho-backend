package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/sonyamoonglade/sancho-backend/auth"
	"github.com/sonyamoonglade/sancho-backend/internal/domain"
)

const (
	HeaderAuthorization = "Authorization"

	ResponseUnauthorized = "unauthorized"
	ResponseAccessDenied = "access denied"
)

// JWTAuthMiddleware validates incoming Bearer access token.
// It also does RBAC, checking incoming user role for required
// and enriches request's context with userID of the token owner.
type JWTAuthMiddleware struct {
	tokenProvider auth.TokenProvider
}

func NewJWTAuthMiddleware(tokenProvider auth.TokenProvider) *JWTAuthMiddleware {
	return &JWTAuthMiddleware{tokenProvider: tokenProvider}
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
			case errors.Is(err, auth.ErrTokenExpired),
				errors.Is(err, auth.ErrInvalidToken),
				errors.Is(err, auth.ErrInvalidIssuer):
				return c.Status(http.StatusUnauthorized).SendString(ResponseUnauthorized)
			}
			// Delegate internal errors to ErrorHandler
			return err
		}

		// Role's insufficient
		if !userAuth.Role.CheckPermissions(requiredRole) {
			return c.Status(http.StatusForbidden).SendString(ResponseAccessDenied)
		}

		return c.Next()
	}
}
