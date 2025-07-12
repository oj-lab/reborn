package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/oj-lab/go-webmods/app"
	config "github.com/oj-lab/reborn/configs"
	"github.com/oj-lab/reborn/internal/middlewares"
	"github.com/oj-lab/reborn/internal/routers"
	"github.com/oj-lab/reborn/internal/services"
)

func init() {
	app.SetCMDName("web")
	cwd, _ := os.Getwd()
	app.Init(cwd)
}

func main() {
	cfg := config.Load()

	// Initialize service manager
	serviceManager := services.NewServiceManager()
	if err := serviceManager.Initialize(cfg); err != nil {
		fmt.Printf("Failed to initialize services: %v\n", err)
		os.Exit(1)
	}

	defer func() {
		if err := serviceManager.Shutdown(); err != nil {
			fmt.Printf("Error during service shutdown: %v\n", err)
		}
	}()

	e := echo.New()

	// Set custom error handler
	e.HTTPErrorHandler = middlewares.ErrorHandler

	// Add middlewares
	e.Use(middlewares.RequestID())
	e.Use(middlewares.Logger())
	e.Use(middlewares.Recover())
	e.Use(middlewares.CORS())
	e.Use(middlewares.RateLimiter())

	// Add health check endpoint
	e.GET("/health", func(c echo.Context) error {
		health := serviceManager.HealthCheck()
		return c.JSON(200, map[string]any{
			"status":   "ok",
			"services": health,
		})
	})

	// Register API routes
	routers.RegisterAPIv1Routes(e, serviceManager)
	routers.RegisterAuthRoutes(e, serviceManager)
	routers.RegisterAdminRoutes(e, serviceManager)

	// Add static website middleware (this should be last)
	e.Use(middlewares.StaticWebsite(cfg))

	// Start server in a goroutine
	go func() {
		if err := e.Start(fmt.Sprintf(":%d", cfg.Server.Port)); err != nil {
			fmt.Printf("Server startup error: %v\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	fmt.Println("Shutting down server...")

	// Shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		fmt.Printf("Server shutdown error: %v\n", err)
	}

	fmt.Println("Server stopped")
}
