package main

import (
	"fmt"
	"os"
	"path"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/oj-lab/reborn/cmd/web/routers"
	"github.com/oj-lab/reborn/common/app"
)

const configKeyServerPort = "server.port"

var port uint

func init() {
	cwd, _ := os.Getwd()
	app.InitConfig(path.Join(cwd, "configs", "web"))
	port = app.Config().GetUint(configKeyServerPort)
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	routers.RegisterAPIv1Routes(e)
	routers.RegisterAuthRoutes(e)

	e.Start(fmt.Sprintf(":%d", port))
}
