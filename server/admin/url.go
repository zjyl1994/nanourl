package admin

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zjyl1994/nanourl/model/val_obj"
	"github.com/zjyl1994/nanourl/service"
	"github.com/zjyl1994/nanourl/vars"
)

type createUrlReqForm struct {
	LongUrl   string `form:"long_url"`
	ShortCode string `form:"short_code"`
}

func CreateUrlHandler(c *fiber.Ctx) error {
	var req createUrlReqForm
	if err := c.BodyParser(&req); err != nil {
		return err
	}
	if len(req.LongUrl) == 0 {
		return fiber.ErrBadRequest
	}

	var svc service.URLService
	id, err := svc.New(val_obj.URLObject{
		LongURL:   req.LongUrl,
		ShortCode: req.ShortCode,
	})
	if err != nil {
		return err
	}
	if len(req.ShortCode) != 0 {
		return c.SendString(vars.BaseUrl + req.ShortCode)
	}
	hashCode, err := svc.GetShortCode(id)
	if err != nil {
		return err
	}
	return c.SendString(vars.BaseUrl + hashCode)
}