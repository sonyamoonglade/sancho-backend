package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/sonyamoonglade/sancho-backend/internal/handler/input"
	"github.com/sonyamoonglade/sancho-backend/internal/validation"
)

func (h Handler) CreateUserOrder(c *fiber.Ctx) error {
	var inp input.CreateUserOrderInput
	if err := c.BodyParser(&inp); err != nil {
		return err
	}
	if ok, msg := validation.ValidateStruct(inp); !ok {
		return c.Status(http.StatusBadRequest).SendString(msg)
	}
	if ok, msg := validation.ValidatePayType(inp.Pay); !ok {
		return c.Status(http.StatusBadRequest).SendString(msg)
	}

	//h.services.Order.CreateUserOrder()

	return nil
}
