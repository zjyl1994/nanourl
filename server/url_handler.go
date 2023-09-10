package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zjyl1994/nanourl/model/val_obj"
	"github.com/zjyl1994/nanourl/service"
	"github.com/zjyl1994/nanourl/vars"
)

func RedirectHandler(c *fiber.Ctx) error {
	shortCode := c.Params("code")
	if shortCode == "" {
		return fiber.ErrBadRequest
	}

	var urlSvc service.URLService
	obj, err := urlSvc.SearchCode(shortCode)
	if err != nil {
		if err == service.ErrCodeNotFound {
			return c.Status(fiber.StatusNotFound).SendString("url not found or expired")
		}
		return err
	}

	var logSvc service.LogService
	var realIp string
	if len(vars.RealIpHeader) > 0 {
		realIp = string(c.Request().Header.Peek(vars.RealIpHeader))
		if realIp == "" {
			realIp = c.Context().RemoteIP().String()
		}
	} else {
		realIp = c.Context().RemoteIP().String()
	}

	logSvc.AddLog(val_obj.AccessLog{
		UrlId:     obj.Id,
		Referrer:  string(c.Request().Header.Referer()),
		UserAgent: string(c.Request().Header.UserAgent()),
		UserIp:    realIp,
	})

	return c.Redirect(obj.LongURL)
}
