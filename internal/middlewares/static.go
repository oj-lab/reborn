package middlewares

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
	config "github.com/oj-lab/reborn/configs"
)

// StaticWebsite serves the frontend application
func StaticWebsite(cfg config.Config) echo.MiddlewareFunc {
	return ServeStaticFiles(cfg.Website.DistPath)
}

// ServeStaticFiles serves static files from the dist directory
func ServeStaticFiles(distPath string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			path := c.Request().URL.Path

			// Skip for API routes and auth routes
			if strings.HasPrefix(path, "/api/") ||
				strings.HasPrefix(path, "/auth/") ||
				strings.HasPrefix(path, "/health") {
				return next(c)
			}

			// Try to serve static file
			filePath := filepath.Join(distPath, path)

			// Check if file exists
			if info, err := os.Stat(filePath); err == nil && !info.IsDir() {
				return serveFile(c, filePath)
			}

			// If it's a directory, try index.html
			if info, err := os.Stat(filePath); err == nil && info.IsDir() {
				indexPath := filepath.Join(filePath, "index.html")
				if _, err := os.Stat(indexPath); err == nil {
					return serveFile(c, indexPath)
				}
			}

			// For SPA routing, serve index.html for unknown routes
			indexPath := filepath.Join(distPath, "index.html")
			if _, err := os.Stat(indexPath); err == nil {
				return serveFile(c, indexPath)
			}

			return next(c)
		}
	}
}

// serveFile serves a file with appropriate content type
func serveFile(c echo.Context, filePath string) error {
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
