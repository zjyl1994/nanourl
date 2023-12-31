package server

import (
	_ "embed"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
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
		ViewsLayout:  "layout",
		ServerHeader: "Nanourl",
	})
	app.Use(favicon.New(favicon.Config{Data: faviconData}))

	adminG := app.Group("/admin")
	if len(vars.AdminUsername) > 0 || len(vars.AdminPassword) > 0 { // add auth when username/password setuped
		adminG.Use(basicauth.New(basicauth.Config{Users: map[string]string{vars.AdminUsername: vars.AdminPassword}}))
	}
	adminG.Static("/static", "./static")
	adminG.Post("/create", admin.CreateUrlHandler)
	adminG.Get("/new", showPage("new"))
	adminG.Get("/log", showPage("log"))
	adminG.Get("/url", showPage("list"))
	adminG.Get("/url/edit", showPage("edit"))
	adminG.Get("/stat", admin.StatisticsPage)
	adminG.Get("/qrcode", admin.GenQRCodeHandler)

	apiG := adminG.Group("/api")
	apiG.Get("/log", admin.ListLogHandler)
	apiG.Get("/url", admin.ListUrlHandler)
	apiG.Post("/toggle", admin.ToggleUrlHandler)
	apiG.Delete("/url", admin.DeleteUrlHandler)
	apiG.Put("/url", admin.UpdateUrlHandler)
	apiG.Get("/url_info", admin.GetUrlHandler)

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

func showPage(name string, layout ...string) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.Render(name, fiber.Map{}, layout...)
	}
}
