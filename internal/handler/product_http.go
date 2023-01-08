package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/sonyamoonglade/sancho-backend/internal/handler/input"
	"github.com/sonyamoonglade/sancho-backend/internal/validation"
)

func (h Handler) GetCatalog(c *fiber.Ctx) error {
	catalog, err := h.services.Product.GetAll(c.Context())
	if err != nil {
		return err
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"catalog": catalog,
	})
}

func (h Handler) GetCategories(c *fiber.Ctx) error {
	var (
		sortedStr = c.Query("sorted", "1" /* true by default */)
		sorted    = sortedStr == "1"
	)
	categories, err := h.services.Product.GetAllCategories(c.Context(), sorted)
	if err != nil {
		return err
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"categories": categories,
	})
}

func (h Handler) CreateProduct(c *fiber.Ctx) error {
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

func (h Handler) DeleteProduct(c *fiber.Ctx) error {
	productID := c.Params("id", "")
	if productID == "" {
		return c.Status(http.StatusBadRequest).SendString("empty id")
	}
	if err := h.services.Product.Delete(c.Context(), productID); err != nil {
		return err
	}
	return c.SendStatus(http.StatusOK)
}

func (h Handler) UpdateProduct(c *fiber.Ctx) error {

	return nil
}

func (h Handler) ApproveProduct(c *fiber.Ctx) error {
	productID := c.Params("id", "")
	if productID == "" {
		return c.Status(http.StatusBadRequest).SendString("empty id")
	}
	if err := h.services.Product.Approve(c.Context(), productID); err != nil {
		return err
	}
	return c.SendStatus(http.StatusOK)
}

func (h Handler) DisapproveProduct(c *fiber.Ctx) error {
	productID := c.Params("id", "")
	if productID == "" {
		return c.Status(http.StatusBadRequest).SendString("empty id")
	}
	if err := h.services.Product.Disapprove(c.Context(), productID); err != nil {
		return err
	}
	return c.SendStatus(http.StatusOK)
}
