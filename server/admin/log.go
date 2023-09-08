package admin

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zjyl1994/nanourl/service"
	"github.com/zjyl1994/nanourl/util"
)

func ListLogHandler(c *fiber.Ctx) error {
	page, pageSize := util.PageNormalize(c.QueryInt("page"), c.QueryInt("size"))
	var logSvc service.LogService
	logs, total, err := logSvc.List(page, pageSize)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"total": total,
		"logs":  logs,
	})
}
