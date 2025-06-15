package routers

import (
	"github.com/labstack/echo/v4"
	"github.com/oj-lab/reborn/cmd/web/handlers"
)

func RegisterAuthRoutes(e *echo.Echo) {
	baseGroup := e.Group("/auth")
	{
		baseGroup.GET("/:provider/login", handlers.Login)
		baseGroup.GET("/:provider/callback", handlers.Callback)
	}
}
