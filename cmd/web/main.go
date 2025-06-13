package main

import (
	"github.com/gin-gonic/gin"
	"github.com/oj-lab/reborn/cmd/web/handlers"
)

func main() {
	engine := gin.Default()
	// Register routes
	userGroup := engine.Group("/user")
	{
		userGroup.POST("/", handlers.CreateUser)
	}
	engine.GET("/login", handlers.LoginUser)
	engine.GET("/callback", handlers.CallbackUser)
	engine.Run(":8080")
}
