package view

import (
	"embed"

	"github.com/labstack/echo/v4"
)

var (
	//go:embed all:dist
	distDir           embed.FS
	DistDirFileSystem = echo.MustSubFS(distDir, "dist")
)
