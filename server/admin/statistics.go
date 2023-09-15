package admin

import "github.com/gofiber/fiber/v2"

func StatisticsPage(c *fiber.Ctx) error {
	return c.Render("statistics", fiber.Map{})
}
