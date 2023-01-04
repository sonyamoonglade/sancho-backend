package handler

import (
	"github.com/gofiber/fiber/v2"
	service "github.com/sonyamoonglade/sancho-backend/internal/services"
)

type Handler struct {
	services *service.Services
}

func NewHandler(services *service.Services) *Handler {
	return &Handler{
		services: services,
	}
}

func (h Handler) InitAPI(router fiber.Router) {
	api := router.Group("/api")
	{
		h.initProductAPI(api)
	}
}

// TODO: move
func (h Handler) errorHandler(c *fiber.Ctx, err error) error {
	return nil
}
