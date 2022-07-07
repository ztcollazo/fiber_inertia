package main

import (
	"embed"
	"net/http"

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

	app.Get("/:name", func(c *fiber.Ctx) error {
		return c.Render("Index", fiber.Map{
			"name": c.Params("name", "world"),
		})
	})

	app.Listen(":8080")
}
