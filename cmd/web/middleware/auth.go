package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.String(http.StatusUnauthorized, "Missing Authorization header")
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.String(http.StatusUnauthorized, "Invalid Authorization header format")
		}

		token := parts[1]
		// TODO: In a real application, you would validate the token (e.g., JWT) here.
		// For this project, we are using a simple token format: user_{id}_token.
		var userID uint64
		if _, err := fmt.Sscanf(token, "user_%d_token", &userID); err != nil {
			return c.String(http.StatusUnauthorized, "Invalid token")
		}

		if userID == 0 {
			return c.String(http.StatusUnauthorized, "Invalid user ID in token")
		}

		c.Set("userID", userID)
		return next(c)
	}
}
