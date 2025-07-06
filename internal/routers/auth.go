package routers

import (
	"github.com/labstack/echo/v4"
	"github.com/oj-lab/reborn/internal/handlers"
	"github.com/oj-lab/reborn/internal/services"
)

// RegisterAuthRoutes registers authentication related routes
func RegisterAuthRoutes(e *echo.Echo, serviceManager *services.ServiceManager) {
	authHandler := handlers.NewAuthHandler(serviceManager.GetAuthService())

	authGroup := e.Group("/auth")
	{
		authGroup.GET("/login", authHandler.Login)
		authGroup.GET("/callback", authHandler.Callback)
		authGroup.POST("/logout", authHandler.Logout)
		authGroup.GET("/logout", authHandler.Logout) // Support both GET and POST for logout
	}
}
