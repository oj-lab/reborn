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
	engine.Run(":8080")
}
