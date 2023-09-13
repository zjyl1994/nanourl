package admin

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zjyl1994/nanourl/service"
	"github.com/zjyl1994/nanourl/util"
)

func ListLogPage(c *fiber.Ctx) error {
	page, pageSize := util.PageNormalize(c.QueryInt("page"), c.QueryInt("size"))
	var logSvc service.LogService
	logs, total, err := logSvc.List(c.QueryInt("id"), page, pageSize)
	if err != nil {
		return err
	}
	return c.Render("log", fiber.Map{"list": logs, "total": total})
}
