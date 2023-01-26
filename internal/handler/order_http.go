package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/sonyamoonglade/sancho-backend/internal/handler/input"
	"github.com/sonyamoonglade/sancho-backend/internal/handler/middleware"
	"github.com/sonyamoonglade/sancho-backend/internal/validation"
	"github.com/sonyamoonglade/sancho-backend/pkg/logger"
	"go.uber.org/zap"
)

func (h Handler) CreateUserOrder(c *fiber.Ctx) error {
	logger.Get().Debug("", zap.String("http order", "here!"))
	var inp input.CreateUserOrderInput
	if err := c.BodyParser(&inp); err != nil {
		return err
	}
	if ok, msg := validation.ValidateStruct(inp); !ok {
		return c.Status(http.StatusBadRequest).SendString(msg)
	}
	if ok, msg := validation.ValidatePayMethod(inp.Pay); !ok {
		return c.Status(http.StatusBadRequest).SendString(msg)
	}
	if ok, msg := validation.ValidateCart(inp.Cart); !ok {
		return c.Status(http.StatusBadRequest).SendString(msg)
	}

	customerID, err := middleware.GetUserIDFromCtx(c)
	if err != nil {
		return err
	}
	orderID, err := h.services.Order.CreateUserOrder(c.Context(), inp.ToDTO(customerID))
	if err != nil {
		return err
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"orderId": orderID,
	})
}
