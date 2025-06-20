package main

import (
	"fmt"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/oj-lab/reborn/cmd/web/routers"
	"github.com/oj-lab/reborn/common/app"
)

const configKeyServerPort = "server.port"

var port uint

func main() {
	cwd, _ := os.Getwd()
	app.Init(cwd, "web")
	port = app.Config().GetUint(configKeyServerPort)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	routers.RegisterAPIv1Routes(e)
	routers.RegisterAuthRoutes(e)

	e.Start(fmt.Sprintf(":%d", port))
}
