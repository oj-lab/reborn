package routers

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
	config "github.com/oj-lab/reborn/configs"
	"github.com/oj-lab/reborn/internal/middlewares"
	"github.com/oj-lab/reborn/internal/services"
)

// RegisterPageRoutes registers all page routes including:
// - Home page (/) - no authentication required
// - Admin pages (/admin/*) - admin authentication required
// - Static file serving for other routes (assets, etc.)
func RegisterPageRoutes(e *echo.Echo, serviceManager *services.ServiceManager) {
	authService := serviceManager.GetAuthService()
	cfg := config.Load()

	// Register home page routes (no authentication required)
	homeHandler := func(c echo.Context) error {
		// Serve the index.html file for home page
		indexPath := filepath.Join(cfg.Website.DistPath, "index.html")
		return serveIndexFile(c, indexPath)
	}

	// Home page routes
	e.GET("/", homeHandler)

	// Register admin page routes with authentication
	adminPageGroup := e.Group("/admin")
	adminPageGroup.Use(middlewares.LoginSession(authService))

	// Admin route handler that serves the frontend index.html
	adminHandler := func(c echo.Context) error {
		// Serve the index.html file for all admin routes
		// The frontend router will handle the specific admin pages
		indexPath := filepath.Join(cfg.Website.DistPath, "index.html")
		return serveIndexFile(c, indexPath)
	}

	// Custom middleware to handle admin authentication for page routes
	adminPageGroup.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Check if user is authenticated
			if !middlewares.IsAuthenticated(c) {
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

	// Admin page routes
	adminPageGroup.GET("", func(c echo.Context) error {
		return c.Redirect(http.StatusFound, "/admin/")
	})

	adminPageGroup.GET("/", adminHandler)
	adminPageGroup.GET("/*", adminHandler)

	// Register static file serving middleware for other routes
	e.Use(serveStaticFiles(cfg.Website.DistPath))
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

// serveStaticFiles serves static files from the dist directory (similar to StaticWebsite middleware)
func serveStaticFiles(distPath string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			path := c.Request().URL.Path

			// Skip for API routes, auth routes, admin routes (handled above), and health
			if strings.HasPrefix(path, "/api/") ||
				strings.HasPrefix(path, "/auth/") ||
				strings.HasPrefix(path, "/admin/") ||
				strings.HasPrefix(path, "/health") ||
				path == "/" { // Home page is handled above
				return next(c)
			}

			// Try to serve static file
			filePath := filepath.Join(distPath, path)

			// Check if file exists
			if info, err := os.Stat(filePath); err == nil && !info.IsDir() {
				return serveStaticFile(c, filePath)
			}

			// If it's a directory, try index.html
			if info, err := os.Stat(filePath); err == nil && info.IsDir() {
				indexPath := filepath.Join(filePath, "index.html")
				if _, err := os.Stat(indexPath); err == nil {
					return serveStaticFile(c, indexPath)
				}
			}

			// For SPA routing, serve index.html for unknown routes
			indexPath := filepath.Join(distPath, "index.html")
			if _, err := os.Stat(indexPath); err == nil {
				return serveIndexFile(c, indexPath)
			}

			return next(c)
		}
	}
}

// serveStaticFile serves a static file with appropriate content type
func serveStaticFile(c echo.Context, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Set content type based on file extension
	ext := filepath.Ext(filePath)
	switch ext {
	case ".html":
		c.Response().Header().Set("Content-Type", "text/html; charset=utf-8")
	case ".css":
		c.Response().Header().Set("Content-Type", "text/css; charset=utf-8")
	case ".js":
		c.Response().Header().Set("Content-Type", "application/javascript; charset=utf-8")
	case ".json":
		c.Response().Header().Set("Content-Type", "application/json; charset=utf-8")
	case ".png":
		c.Response().Header().Set("Content-Type", "image/png")
	case ".jpg", ".jpeg":
		c.Response().Header().Set("Content-Type", "image/jpeg")
	case ".gif":
		c.Response().Header().Set("Content-Type", "image/gif")
	case ".svg":
		c.Response().Header().Set("Content-Type", "image/svg+xml")
	case ".ico":
		c.Response().Header().Set("Content-Type", "image/x-icon")
	default:
		c.Response().Header().Set("Content-Type", "application/octet-stream")
	}

	// Set cache headers for static assets
	if strings.Contains(filePath, "/assets/") {
		c.Response().Header().Set("Cache-Control", "public, max-age=31536000") // 1 year
	} else {
		c.Response().Header().Set("Cache-Control", "no-cache")
	}

	c.Response().WriteHeader(http.StatusOK)
	_, err = io.Copy(c.Response().Writer, file)
	return err
}
