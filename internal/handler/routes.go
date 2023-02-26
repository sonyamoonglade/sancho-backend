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
		products.Post("/create", h.AdminCreateProduct)
		products.Put("/:id/update", h.AdminUpdateProduct)
		products.Delete("/:id/delete", h.AdminDeleteProduct)
		products.Put("/:id/approve", h.AdminApproveProduct)
		products.Put("/:id/disapprove", h.AdminDisapproveProduct)
	}
}

func (h Handler) initOrdersAPI(api fiber.Router) {
	m := h.middlewares

	customerAuth := m.JWTAuth.Use(domain.RoleCustomer)
	// Don't order.Use(customerAuth), because not all /order requests require auth.
	order := api.Group("/order")
	{
		order.Post("/create", customerAuth, h.CreateUserOrder)

		{
			worker := order.Group("/worker")
			worker.Use(m.JWTAuth.Use(domain.RoleWorker))

			worker.Post("/create", h.CreateWorkerOrder)
		}
	}
}
