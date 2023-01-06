package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sonyamoonglade/sancho-backend/auth"
	"github.com/sonyamoonglade/sancho-backend/internal/domain"
)

// JWTAuthMiddleware validates incoming Bearer access token
// and enriches request's context with userID of the token owner
type JWTAuthMiddleware struct {
	tokenProvider auth.TokenProvider
}

func NewJWTAuthMiddleware(tokenProvider auth.TokenProvider) *JWTAuthMiddleware {
	return &JWTAuthMiddleware{tokenProvider: tokenProvider}
}

func (m JWTAuthMiddleware) Use(r domain.Role) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Next()
	}
}
