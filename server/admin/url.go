package admin

import (
	"database/sql"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"github.com/zjyl1994/nanourl/model/val_obj"
	"github.com/zjyl1994/nanourl/service"
	"github.com/zjyl1994/nanourl/util"
	"github.com/zjyl1994/nanourl/vars"
)

type createUrlReqForm struct {
	LongUrl      string `form:"long_url"`
	ShortCode    string `form:"short_code"`
	ExpireSecond int    `form:"expire_sec"`
}

func CreateUrlHandler(c *fiber.Ctx) error {
	var req createUrlReqForm
	if err := c.BodyParser(&req); err != nil {
		return err
	}
	if len(req.LongUrl) == 0 {
		return fiber.ErrBadRequest
	}
	valobj := val_obj.URLObject{
		LongURL:   req.LongUrl,
		ShortCode: req.ShortCode,
	}
	if req.ExpireSecond > 0 {
		valobj.ExpireTime = sql.NullTime{Time: time.Now().Add(time.Duration(req.ExpireSecond) * time.Second), Valid: true}
	}

	var svc service.URLService
	_, shortCode, err := svc.New(valobj)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).SendString(vars.BaseUrl + shortCode)
}

func ListUrlPage(c *fiber.Ctx) error {
	page, pageSize := util.PageNormalize(c.QueryInt("page"), c.QueryInt("size"))
	var svc service.URLService
	data, total, err := svc.List(page, pageSize)
	if err != nil {
		return err
	}
	var logsvc service.LogService
	logCount, err := logsvc.CountLog(lo.Map(data, func(x val_obj.URLObject, _ int) int {
		return int(x.Id)
	}))
	if err != nil {
		return err
	}
	return c.Render("admin/list", fiber.Map{"list": data, "total": total, "base_url": vars.BaseUrl, "logc": logCount})
}
