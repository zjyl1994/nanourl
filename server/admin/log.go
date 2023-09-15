package admin

import (
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"github.com/zjyl1994/nanourl/model/render_obj"
	"github.com/zjyl1994/nanourl/model/val_obj"
	"github.com/zjyl1994/nanourl/service"
	"github.com/zjyl1994/nanourl/util"
	"github.com/zjyl1994/nanourl/vars"
)

func ListLogHandler(c *fiber.Ctx) error {
	page, pageSize := util.PageNormalize(c.QueryInt("page"), c.QueryInt("size"))
	var logSvc service.LogService
	logs, total, err := logSvc.List(c.QueryInt("id"), page, pageSize)
	if err != nil {
		return err
	}
	list := lo.Map(logs, func(x val_obj.AccessLog, _ int) render_obj.AccessLog {
		var countryName, countryEmoji string
		if len(x.UserCountry) > 0 {
			if ci, ok := vars.GeoCountry[x.UserCountry]; ok {
				countryName = ci.Name
				countryEmoji = ci.Emoji
			}
		}

		return render_obj.AccessLog{
			UrlId:        x.UrlId,
			Referrer:     x.Referrer,
			UserIp:       x.UserIp,
			CountryCode:  x.UserCountry,
			CountryName:  countryName,
			CountryEmoji: countryEmoji,
			UserAgent:    x.UserAgent,
			AccessTime:   x.AccessTime.Unix(),
			OS:           x.OS,
			Browser:      x.Browser,
			Device:       x.Device,
			DeviceType:   x.DeviceType,
		}
	})
	return c.JSON(fiber.Map{
		"list":       list,
		"page":       page,
		"page_size":  pageSize,
		"total_rows": total,
	})
}
