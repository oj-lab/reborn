package routers

import (
	"github.com/labstack/echo/v4"
	"github.com/oj-lab/reborn/cmd/web/handlers"
	"github.com/oj-lab/reborn/cmd/web/middleware"
)

func RegisterAuthRoutes(e *echo.Echo) {
	baseGroup := e.Group("/auth")
	{
		// TODO: Remove `provider` here, use OAuth state and store provider type in it
		baseGroup.GET("/:provider/login", handlers.Login)
		baseGroup.GET("/:provider/callback", handlers.Callback)

		baseGroup.POST("/login", handlers.LoginWithPassword)
		baseGroup.POST("/register", handlers.RegisterWithPassword)

		authRequired := baseGroup.Group("")
		authRequired.Use(middleware.Auth)
		{
			authRequired.POST("/password", handlers.SetPassword)
		}
	}
}
