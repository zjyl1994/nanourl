package server

import "github.com/gofiber/fiber/v2"

func Run(listen string) error {
	app := fiber.New()
	
	return app.Listen(listen)
}
