package admin

import (
	"database/sql"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	qrcode "github.com/skip2/go-qrcode"
	"github.com/zjyl1994/nanourl/model/render_obj"
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

func ListUrlHandler(c *fiber.Ctx) error {
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
	list := make([]render_obj.URLObject, len(data))
	for i, v := range data {
		list[i] = render_obj.URLObject{
			Id:         v.Id,
			LongURL:    v.LongURL,
			ShortCode:  v.ShortCode,
			CreateTime: v.CreateTime.Unix(),
			ExpireTime: lo.Ternary(v.ExpireTime.Valid, v.ExpireTime.Time.Unix(), 0),
			Enabled:    v.Enabled,
			ClickCount: lo.ValueOr(logCount, v.Id, 0),
			HrefLink:   vars.BaseUrl + v.ShortCode,
		}
	}
	return c.JSON(fiber.Map{
		"list":       list,
		"page":       page,
		"page_size":  pageSize,
		"total_rows": total,
	})
}

func GenQRCodeHandler(c *fiber.Ctx) error {
	shortCode := c.Query("code")
	if shortCode == "" {
		return fiber.ErrBadRequest
	}
	targetUrl := vars.BaseUrl + shortCode
	png, err := qrcode.Encode(targetUrl, qrcode.Medium, 256)
	if err != nil {
		return err
	}
	c.Set(fiber.HeaderContentType, "image/png")
	return c.Send(png)
}

type updateCodeEnabledReq struct {
	Id      uint `form:"id"`
	Enabled bool `form:"enabled"`
}

func ToggleUrlHandler(c *fiber.Ctx) error {
	var req updateCodeEnabledReq
	if err := c.BodyParser(&req); err != nil {
		return err
	}
	if req.Id == 0 {
		return fiber.ErrBadRequest
	}
	var svc service.URLService
	err := svc.UpdateEnabled(req.Id, req.Enabled)
	if err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusOK)
}

func DeleteUrlHandler(c *fiber.Ctx) error {
	id := c.QueryInt("id")
	if id <= 0 {
		return fiber.ErrBadRequest
	}
	var svc service.URLService
	err := svc.Delete(uint(id))
	if err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusOK)
}
