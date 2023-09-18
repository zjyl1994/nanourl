package admin

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zjyl1994/nanourl/service"
)

func StatisticsPage(c *fiber.Ctx) error {
	var svc service.StatService
	topUrls, err := svc.TopURL()
	if err != nil {
		return err
	}
	topCountrys, err := svc.TopCountry()
	if err != nil {
		return err
	}
	return c.Render("statistics", fiber.Map{
		"top_url":     topUrls,
		"top_country": topCountrys,
	})
}
