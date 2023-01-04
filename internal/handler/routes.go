package handler

import "github.com/gofiber/fiber/v2"

func (h Handler) initProductAPI(api fiber.Router) {
	p := api.Group("/products")
	{
		p.Get("/catalog")
		p.Get("/categories")

		adm := p.Group("/a")
		{
			adm.Post("/create")
			adm.Delete("/:id/delete")
			adm.Put("/:id/update")
			adm.Put("/:id/approve")
			adm.Put("/:id/changeImageUrl")
		}
	}

}
