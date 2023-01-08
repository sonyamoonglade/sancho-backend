package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sonyamoonglade/sancho-backend/internal/domain"
)

func (h Handler) initProductAPI(api fiber.Router) {
	m := h.middlewares

	p := api.Group("/products")
	{
		p.Get("/catalog", h.GetCatalog)
		p.Get("/categories", h.GetCategories)

		adm := p.Group("/a")
		{
			adm.Use(m.JWTAuth.Use(domain.RoleAdmin))

			adm.Post("/create", h.CreateProduct)
			adm.Delete("/:id/delete", h.DeleteProduct)
			adm.Put("/:id/update", h.UpdateProduct)
			adm.Put("/:id/approve", h.ApproveProduct)
			adm.Put("/:id/disapprove", h.DisapproveProduct)
		}
	}

}
