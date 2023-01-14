package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/sonyamoonglade/sancho-backend/internal/handler/input"
	"github.com/sonyamoonglade/sancho-backend/internal/handler/middleware"
	"github.com/sonyamoonglade/sancho-backend/internal/validation"
)

func (h Handler) AdminCreateProduct(c *fiber.Ctx) error {
	var inp input.CreateProductInput
	if err := c.BodyParser(&inp); err != nil {
		return err
	}

	if ok, msg := validation.ValidateFeatures(inp.Features); !ok {
		return c.Status(http.StatusBadRequest).SendString(msg)
	}

	productID, err := h.services.Product.Create(c.Context(), inp.ToDTO())
	if err != nil {
		return err
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"productId": productID,
	})
}

func (h Handler) AdminDeleteProduct(c *fiber.Ctx) error {
	productID := c.Params("id", "")
	if productID == "" {
		return c.Status(http.StatusBadRequest).SendString("empty id")
	}
	if err := h.services.Product.Delete(c.Context(), productID); err != nil {
		return err
	}
	return c.SendStatus(http.StatusOK)
}

func (h Handler) AdminUpdateProduct(c *fiber.Ctx) error {
	productID := c.Params("id", "")
	if productID == "" {
		return c.Status(http.StatusBadRequest).SendString("empty id")
	}
	var inp input.UpdateProductInput
	if err := c.BodyParser(&inp); err != nil {
		return err
	}
	if err := h.services.Product.Update(c.Context(), inp.ToDTO(productID)); err != nil {
		return err
	}
	return c.SendStatus(http.StatusOK)
}

func (h Handler) AdminApproveProduct(c *fiber.Ctx) error {
	productID := c.Params("id", "")
	if productID == "" {
		return c.Status(http.StatusBadRequest).SendString("empty id")
	}
	if err := h.services.Product.Approve(c.Context(), productID); err != nil {
		return err
	}
	return c.SendStatus(http.StatusOK)
}

func (h Handler) AdminDisapproveProduct(c *fiber.Ctx) error {
	productID := c.Params("id", "")
	if productID == "" {
		return c.Status(http.StatusBadRequest).SendString("empty id")
	}
	if err := h.services.Product.Disapprove(c.Context(), productID); err != nil {
		return err
	}
	return c.SendStatus(http.StatusOK)
}

func (h Handler) AdminRefresh(c *fiber.Ctx) error {
	adminID, err := middleware.GetUserIDFromCtx(c)
	if err != nil {
		return err
	}
	_ = adminID
	return nil
}
