package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
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
