package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type XRequestIDMiddleware struct{}

func (m XRequestIDMiddleware) Use() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("X-Request-Id", uuid.NewString())
		return c.Next()
	}
}
