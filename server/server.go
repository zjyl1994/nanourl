package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zjyl1994/nanourl/server/admin"
)

func Run(listen string) error {
	app := fiber.New()

	adminG := app.Group("/admin")
	adminG.Post("/create", admin.CreateUrlHandler)

	app.Get("/:code", RedirectHandler)

	return app.Listen(listen)
}
