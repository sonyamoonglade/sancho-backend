package middleware

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

const userIDCtx = "userid"

var (
	ErrInvalidUserIDFormat = errors.New("invalid user id format")
)

type Middlewares struct {
	JWTAuth    *JWTAuthMiddleware
	XRequestID *XRequestIDMiddleware
}

func NewMiddlewares(jwtAuth *JWTAuthMiddleware, xRequestID *XRequestIDMiddleware) *Middlewares {
	return &Middlewares{
		JWTAuth:    jwtAuth,
		XRequestID: xRequestID,
	}
}

func GetUserIDFromCtx(c *fiber.Ctx) (string, error) {
	userID := c.Locals(userIDCtx)
	if _, ok := userID.(string); !ok {
		return "", ErrInvalidUserIDFormat
	}

	return userID.(string), nil
}
