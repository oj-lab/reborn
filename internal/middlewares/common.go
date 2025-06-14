package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// CORS returns a CORS middleware with common settings
func CORS() echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"}, // Configure based on your needs
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
		},
	})
}

// RequestID returns a request ID middleware
func RequestID() echo.MiddlewareFunc {
	return middleware.RequestID()
}

// Logger returns a logger middleware
func Logger() echo.MiddlewareFunc {
	return middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} ${status} ${method} ${uri} ${latency_human}\n",
	})
}

// Recover returns a recover middleware
func Recover() echo.MiddlewareFunc {
	return middleware.Recover()
}

// RateLimiter returns a rate limiter middleware
func RateLimiter() echo.MiddlewareFunc {
	return middleware.RateLimiter(
		middleware.NewRateLimiterMemoryStore(20),
	) // 20 requests per second
}
