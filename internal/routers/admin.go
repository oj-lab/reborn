package routers

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
	config "github.com/oj-lab/reborn/configs"
	"github.com/oj-lab/reborn/internal/middlewares"
	"github.com/oj-lab/reborn/internal/services"
)

// RegisterAdminRoutes registers admin page routes with proper authentication
func RegisterAdminRoutes(e *echo.Echo, serviceManager *services.ServiceManager) {
	authService := serviceManager.GetAuthService()
	cfg := config.Load()

	// Create admin group with middleware
	adminGroup := e.Group("/admin")
	adminGroup.Use(middlewares.LoginSession(authService))

	// Custom middleware to handle admin authentication for page routes
	adminGroup.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Check if user is authenticated
			if !middlewares.IsAuthenticated(c) {
				// Redirect to login instead of returning error for page requests
				return c.Redirect(http.StatusFound, "/auth/login")
			}

			// Get user token and check if user is admin
			userToken := middlewares.GetUserToken(c)
			if userToken == "" {
				return c.Redirect(http.StatusFound, "/auth/login")
			}

			// Use AdminOnly middleware logic but handle errors gracefully
			adminOnlyMiddleware := middlewares.AdminOnly(authService)
			adminHandler := adminOnlyMiddleware(func(c echo.Context) error {
				return nil // Success, user is admin
			})

			if err := adminHandler(c); err != nil {
				// For admin access denied, redirect to home page
				if he, ok := err.(*echo.HTTPError); ok && he.Code == http.StatusForbidden {
					return c.Redirect(http.StatusFound, "/?error=access_denied")
				}
				// For other errors (like service unavailable), redirect to login
				return c.Redirect(http.StatusFound, "/auth/login")
			}

			// User is authenticated and is admin, continue
			return next(c)
		}
	})

	// Admin route handler that serves the frontend index.html
	adminHandler := func(c echo.Context) error {
		// Serve the index.html file for all admin routes
		// The frontend router will handle the specific admin pages
		indexPath := filepath.Join(cfg.Website.DistPath, "index.html")
		return serveIndexFile(c, indexPath)
	}

	// Admin routes
	adminGroup.GET("", func(c echo.Context) error {
		return c.Redirect(http.StatusFound, "/admin/")
	})

	adminGroup.GET("/", adminHandler)
	adminGroup.GET("/*", adminHandler)
}

// serveIndexFile serves the index.html file for SPA routing
func serveIndexFile(c echo.Context, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Admin page not found")
	}
	defer file.Close()

	// Set content type for HTML
	c.Response().Header().Set("Content-Type", "text/html; charset=utf-8")
	c.Response().Header().Set("Cache-Control", "no-cache")

	c.Response().WriteHeader(http.StatusOK)
	_, err = io.Copy(c.Response().Writer, file)
	return err
}
