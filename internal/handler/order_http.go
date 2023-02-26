package handler

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/sonyamoonglade/sancho-backend/internal/domain"
	"github.com/sonyamoonglade/sancho-backend/internal/handler/input"
	"github.com/sonyamoonglade/sancho-backend/internal/handler/middleware"
	"github.com/sonyamoonglade/sancho-backend/internal/services/dto"
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

func (h Handler) CreateWorkerOrder(c *fiber.Ctx) error {
	var inp input.CreateWorkerOrderInput
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

	var customerID string
	customer, err := h.services.User.GetCustomerByPhoneNumber(c.Context(), inp.PhoneNumber)
	if err != nil {
		// If customer not found then should register a new one.
		if errors.Is(err, domain.ErrCustomerNotFound) {
			// begin register
			registerCustomerDTO := dto.RegisterCustomerDTO{
				PhoneNumber:  inp.PhoneNumber,
				CustomerName: &inp.CustomerName,
			}
			// If order is delivered can fulfill user delivery details
			if inp.DeliveryAddress != nil && inp.IsDelivered {
				registerCustomerDTO.DeliveryAddress = inp.DeliveryAddress.ToUserDeliveryAddress()
			}

			customerID, err = h.services.Auth.RegisterCustomer(c.Context(), registerCustomerDTO)
			if err != nil {
				return err
			}
			// end register customer
		} else {
			return err
		}

	} else {
		// If customer exists then assign it's id
		customerID = customer.UserID.Hex()
	}

	orderID, err := h.services.Order.CreateWorkerOrder(c.Context(), inp.ToDTO(customerID))
	if err != nil {
		return err
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"orderId": orderID,
	})
}
