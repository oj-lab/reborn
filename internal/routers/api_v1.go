package routers

import (
	"github.com/labstack/echo/v4"
	"github.com/oj-lab/reborn/internal/handlers"
	"github.com/oj-lab/reborn/internal/middlewares"
	"github.com/oj-lab/reborn/internal/services"
)

// RegisterAPIv1Routes
//
//	@title		API V1
//	@version	1.0
//	@BasePath	/api/v1
func RegisterAPIv1Routes(e *echo.Echo, serviceManager *services.ServiceManager) {
	authService := serviceManager.GetAuthService()

	// Initialize handlers
	userHandler := handlers.NewUserHandler(authService)

	baseGroup := e.Group("/api/v1")
	{
		userGroup := baseGroup.Group("/user")
		userGroup.Use(middlewares.LoginSession(authService))
		{
			userGroup.GET("/me", userHandler.GetCurrentUser)
		}
	}
}
