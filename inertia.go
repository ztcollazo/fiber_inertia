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
	Template   string
}

type Engine struct {
	*html.Engine
	ctx    *fiber.Ctx
	config Config
	props  map[string]any
	next   map[string]any
	params map[string]any
}

func (e *Engine) Share(name string, value any) {
	e.props[name] = value
}

func (e *Engine) AddProp(name string, value any) {
	e.next[name] = value
}

func (e *Engine) AddParam(name string, value any) {
	e.params[name] = value
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

	return e.display(component, p, w)
}

func (e *Engine) display(component string, props fiber.Map, w io.Writer) error {
	dt := e.props
	for k, v := range e.next {
		dt[k] = v
	}
	for k, v := range props {
		dt[k] = v
	}
	data := map[string]interface{}{
		"component": component,
		"props":     dt,
		"url":       e.ctx.OriginalURL(),
		"version":   e.ctx.Get("X-Inertia-Version", ""),
	}

	e.next = map[string]any{}

	renderJSON, err := strconv.ParseBool(e.ctx.Get("X-Inertia", "false"))

	if err != nil {
		log.Fatal("X-Inertia not parsable")
	}

	if renderJSON && e.ctx.XHR() {
		return jsonResponse(e.ctx, data)
	}

	return htmlResponse(data, w, e.config.Template, e.Engine.Render, e.params)
}

func htmlResponse(data fiber.Map, w io.Writer, template string, renderer func(io.Writer, string, any, ...string) error, params map[string]any) error {
	componentDataByte, err := json.Marshal(data)

	if err != nil {
		panic(err)
	}

	vals := fiber.Map{
		"Page": string(componentDataByte),
	}

	for k, v := range params {
		vals[k] = v
	}

	return renderer(w, template, vals)
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
	if cfg.Template == "" {
		cfg.Template = "index"
	}
	return &Engine{
		Engine: engine,
		config: cfg,
		props:  make(map[string]any),
		params: make(map[string]any),
		next:   make(map[string]any),
	}
}
