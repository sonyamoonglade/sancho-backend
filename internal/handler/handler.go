package handler

import (
	"time"

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

func (h Handler) InitAPI() {
	app := fiber.New(fiber.Config{
		Immutable:    true,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
		ErrorHandler: h.errorHandler,
	})

	api := app.Group("/api")
	{
		h.initProductAPI(api)
	}
}

func (h Handler) errorHandler(c *fiber.Ctx, err error) error {
	return nil
}
