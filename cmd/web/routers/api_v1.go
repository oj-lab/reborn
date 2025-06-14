package routers

import (
	"github.com/labstack/echo/v4"
	"github.com/oj-lab/reborn/cmd/web/handlers"
)

// RegisterAPIv1Routes
//
//	@title		API V1
//	@version	1.0
//	@BasePath	/api/v1
func RegisterAPIv1Routes(e *echo.Echo) {
	baseGroup := e.Group("/api/v1")
	{
		userGroup := baseGroup.Group("/user")
		{
			userGroup.POST("", CreateUser)
		}
	}
}

// CreateUser
//
//	@Summary		Create a new user
//	@Description	Create a new user with the provided details
//	@Tags			User
//	@Router			/user [post]
//	@Accept			json
//	@Produce		json
//	@Param			user	body	userpb.CreateUserRequest	true	"User details"
//	@Success		200
func CreateUser(c echo.Context) error {
	return handlers.CreateUser(c)
}
