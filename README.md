# Fiber Inertia

`fiber_inertia` provides a wrapper around the Inertia.js protocol for use with the [Fiber](https://gofiber.io) framework. It provides a views engine that wraps Fiber's default HTML one.

_Note: There are not any official tests due to the complexity of the point of the project. All features have been tested in the example project. If, however, you do find a bug, please report it in the issues_.

Example use:

```go
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
    FS:         http.FS(fs), // or
    Root:       "."
    AssetsPath: "./src",
  })

  app := fiber.New(fiber.Config{
    Views: engine,
  }) 

  app.Use(engine.Middleware())

  app.Get("/:name", func(c *fiber.Ctx) error {
    return c.Render("Index", fiber.Map{
      "name": c.Params("name", "world"),
    })
  })

  app.Listen(":8080")
}
```

And then in the `index.html`

```hbs
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <link rel="icon" type="image/svg+xml" href="/favicon.ico" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>App</title>
  </head>
  <body>
    <div id="app" data-page="{{.Page}}"></div>
    <!--This is just an example. For a more in-depth one using Vite, see the example folder.-->
    <script type="module" src="/app.js"></script>
  </body>
</html>
```

You can then go about with your typical Inertia client-side setup, just like you would normally. See [here](https://inertiajs.com/client-side-setup) for more information.

## API

Fiber Inertia only exports three things: The engine, the config, and the New function.

### `New(cfg Config) *Engine`

`New` creates a new engine based on the configuration provided. See below for more config options.

### `type Config struct`

`Config` has four properties:

- `Root` - The root directory. Do not use this if you are using FS.
- `FS` - the `http.FileSystem`. You can use `embed.FS` and then call `http.FS` on it.
- `AssetsPath` - The path to the assets to version.
- `Template` - The path to the root HTML template that will be rendered with the page string.

### `type Engine struct`

Engine extends `github.com/gofiber/template/html.Engine` to provide a typical HTML-y experience. It also has a Middleware function that's usage is **required**. See the example for more information. You can pass this to the `Views` property in your Fiber app's config. It provides a single variable for your HTML template, `Page`. You can use it like so:

```hbs
<div id="app" data-page="{{.Page}}"></div>
```

Extra functions:

- `Share`: Shares a prop for every request.
- `AddProp`: Add a prop from middleware for the next request.
- `AddParam`: Share a param with the root template.

## License

This project is licensed under the MIT license. See the License.txt for more details.

## Credits

Most of this project is taken from <https://github.com/theArtechnology/fiber-inertia>, except it has been modified to allow its use as an acual view engine for Fiber. Credits to [@theArtechnology](https://twitter.com/theArtechnology)
