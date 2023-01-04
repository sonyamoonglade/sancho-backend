package handler

import "github.com/gofiber/fiber/v2"

func (h Handler) initProductAPI(api fiber.Router) {
	p := api.Group("/products")
	{
		p.Get("/catalog", h.GetCatalog)
		p.Get("/categories", h.GetCategories)

		adm := p.Group("/a")
		{
			adm.Post("/create", h.CreateProduct)
			adm.Delete("/:id/delete", h.DeleteProduct)
			adm.Put("/:id/update", h.UpdateProduct)
			adm.Put("/:id/approve", h.ApproveProduct)
			adm.Put("/:id/changeImageUrl", h.ChangeImageURL)
		}
	}

}
