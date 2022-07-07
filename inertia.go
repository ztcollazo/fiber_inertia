package fiber_inertia

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
)

type Config struct {
	Root       string
	FS         http.FileSystem
	AssetsPath string
}

type Engine struct {
	*html.Engine
	ctx    *fiber.Ctx
	config Config
}

func (e *Engine) Middleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if e.config.AssetsPath == "" {
			panic("please provide an assets path")
		}

		hash := hashDir(e.config.AssetsPath)

		if c.Method() == "GET" && c.XHR() && c.Get("X-Inertia-Version", "1") != hash {
			c.Set("X-Inertia-Location", c.Path())
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{})
		}

		c.Set("X-Inertia-Version", hash)

		e.ctx = c
		return c.Next()
	}
}

func (e *Engine) Render(w io.Writer, component string, props any, paths ...string) error {
	p := partialReload(e.ctx, component, props.(fiber.Map))

	return display(e.ctx, component, p, w, e.Engine.Render)
}

func display(c *fiber.Ctx, component string, props fiber.Map, w io.Writer, renderer func(io.Writer, string, any, ...string) error) error {
	data := map[string]interface{}{
		"component": component,
		"props":     props,
		"url":       c.OriginalURL(),
		"version":   c.Get("X-Inertia-Version", ""),
	}

	renderJSON, err := strconv.ParseBool(c.Get("X-Inertia", "false"))

	if err != nil {
		log.Fatal("X-Inertia not parsable")
	}

	if renderJSON && c.XHR() {
		return jsonResponse(c, data)
	}

	return htmlResponse(data, w, renderer)
}

func htmlResponse(data fiber.Map, w io.Writer, renderer func(io.Writer, string, any, ...string) error) error {
	componentDataByte, err := json.Marshal(data)

	if err != nil {
		panic(err)
	}

	return renderer(w, "index", fiber.Map{
		"Page": string(componentDataByte),
	})
}

func jsonResponse(c *fiber.Ctx, page fiber.Map) error {
	jsonByte, _ := json.Marshal(page)
	return c.Status(fiber.StatusOK).JSON(string(jsonByte))
}

func partialReload(c *fiber.Ctx, component string, props fiber.Map) fiber.Map {
	if c.Get("X-Inertia-Partial-Component", "/") == component {
		var newProps = make(fiber.Map)
		partials := strings.Split(c.Get("X-Inertia-Partial-Data", ""), ",")
		for key := range props {
			for _, partial := range partials {
				if key == partial {
					newProps[partial] = props[key]
				}
			}
		}
	}
	return props
}

func New(cfg Config) *Engine {
	var engine *html.Engine
	if cfg.FS != nil {
		engine = html.NewFileSystem(cfg.FS, ".html")
	} else {
		engine = html.New(cfg.Root, ".html")
	}
	return &Engine{
		Engine: engine,
		config: cfg,
	}
}
