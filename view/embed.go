package view

import (
	"embed"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

//go:embed all:dist
var distDir embed.FS

func EmbeddedWebApp() echo.MiddlewareFunc {
	fileSystem := http.FS(echo.MustSubFS(distDir, "dist"))

	config := middleware.StaticConfig{
		Skipper: func(c echo.Context) bool {
			return hasPrefix(c.Path(), "api") || hasPrefix(c.Path(), "socket")
		},
		Root:       ".",
		Index:      "index.html",
		HTML5:      true,
		Browse:     false,
		IgnoreBase: false,
		Filesystem: fileSystem,
	}

	return middleware.StaticWithConfig(config)
}

func hasPrefix(path, prefix string) bool {
	return strings.HasPrefix(strings.ToLower(path), strings.ToLower(prefix))
}
