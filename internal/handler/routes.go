package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sonyamoonglade/sancho-backend/internal/domain"
)

func (h Handler) initProductAPI(api fiber.Router) {
	p := api.Group("/products")
	{
		p.Get("/catalog", h.GetCatalog)
		p.Get("/categories", h.GetCategories)
	}
}

func (h Handler) initAdminsAPI(api fiber.Router) {
	m := h.middlewares

	admins := api.Group("/admins")
	admins.Use(m.JWTAuth.Use(domain.RoleAdmin))

	products := admins.Group("/products")
	{
		products.Post("/create", h.CreateProduct)
		products.Put("/:id/update", h.UpdateProduct)
		products.Delete("/:id/delete", h.DeleteProduct)
		products.Put("/:id/approve", h.ApproveProduct)
		products.Put("/:id/disapprove", h.DisapproveProduct)
	}
}
