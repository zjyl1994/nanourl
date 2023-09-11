package server

import (
	_ "embed"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/template/html/v2"
	"github.com/zjyl1994/nanourl/server/admin"
	"github.com/zjyl1994/nanourl/util"
	"github.com/zjyl1994/nanourl/vars"
)

//go:embed favicon.ico
var faviconData []byte

func Run(listen string) error {
	engine := html.New("./views", ".html")
	engine.Reload(true)
	engine.AddFuncMap(map[string]interface{}{
		"time_f":      util.FormatTime,
		"null_time_f": util.FormatNullableTime,
	})

	app := fiber.New(fiber.Config{
		Views:        engine,
		ServerHeader: "Nanourl",
	})
	app.Use(favicon.New(favicon.Config{Data: faviconData}))

	adminG := app.Group("/admin")
	adminG.Post("/create", admin.CreateUrlHandler)
	adminG.Get("/log", admin.ListLogPage)
	adminG.Get("/url", admin.ListUrlPage)
	adminG.Get("/qrcode", admin.GenQRCodeHandler)

	app.Get("/", indexHandler)
	app.Get("/:code", RedirectHandler)

	return app.Listen(listen)
}

func indexHandler(c *fiber.Ctx) error {
	if vars.HomepageUrl != "" {
		return c.Redirect(vars.HomepageUrl)
	}
	return c.Redirect("admin/url")
}
