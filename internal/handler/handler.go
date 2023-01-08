package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sonyamoonglade/sancho-backend/internal/handler/middleware"
	service "github.com/sonyamoonglade/sancho-backend/internal/services"
)

type Handler struct {
	services    *service.Services
	middlewares *middleware.Middlewares
}

func NewHandler(services *service.Services, middlewares *middleware.Middlewares) *Handler {
	return &Handler{
		services:    services,
		middlewares: middlewares,
	}
}

func (h Handler) InitAPI(router fiber.Router) {
	m := h.middlewares
	api := router.Group("/api")
	api.Use(m.XRequestID.Use())
	{
		h.initProductAPI(api)
	}
}

// TODO: move
func (h Handler) errorHandler(c *fiber.Ctx, err error) error {
	return nil
}
