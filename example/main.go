package main

import (
	"embed"
	"net/http"
	"net/url"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	fi "github.com/ztcollazo/fiber_inertia"
)

//go:embed *
var fs embed.FS

func main() {
	engine := fi.New(fi.Config{
		FS:         http.FS(fs),
		AssetsPath: "./src",
	})

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Use(logger.New())
	app.Use(engine.Middleware())

	engine.Share("start", time.Now().String())
	engine.AddParam("Title", "Example App")

	app.Use(func(c *fiber.Ctx) error {
		engine.AddProp("req", time.Now().String())
		return c.Next()
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("Index", fiber.Map{
			"name": "world",
		})
	})

	app.Get("/:name", func(c *fiber.Ctx) error {
		name, err := url.PathUnescape(c.Params("name", "world"))
		if err != nil {
			return err
		}
		return c.Render("Index", fiber.Map{
			"name": name,
		})
	})

	app.Listen(":8080")
}
