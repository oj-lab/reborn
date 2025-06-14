package middlewares

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/oj-lab/reborn/internal/services"
	"github.com/oj-lab/user-service/pkg/userpb"
)

// Context keys for storing user information
const (
	LoginSessionCookieName = "login_session"
	UserTokenKey           = "user_token"
)

// LoginSession returns a middleware that validates login session from cookie
// and stores user token in context for subsequent handlers
func LoginSession(authService *services.AuthService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Check if auth service is available
			if authService == nil || !authService.IsHealthy() {
				// Log the issue for debugging
				c.Logger().Warn("Auth service is not available or healthy")
				// Continue without authentication instead of returning error
				return next(c)
			}

			// Get session cookie
			cookie, err := c.Cookie(LoginSessionCookieName)
			if err != nil {
				// No session cookie found, continue without authentication
				c.Logger().Debug("No session cookie found")
				return next(c)
			}

			sessionID := cookie.Value
			if sessionID == "" {
				// Empty session ID, continue without authentication
				c.Logger().Debug("Empty session ID")
				return next(c)
			}

			// Get user token from auth service using session ID
			client := authService.GetClient()
			if client == nil {
				c.Logger().Warn("Auth client is not available")
				// Continue without authentication instead of returning error
				return next(c)
			}

			userToken, err := client.GetClient().
				GetUserToken(c.Request().Context(), &userpb.GetUserTokenRequest{
					SessionId: sessionID,
				})
			if err != nil {
				// Invalid or expired session, clear cookie and continue
				c.Logger().Debug("Failed to get user token, clearing session cookie")
				c.SetCookie(&http.Cookie{
					Name:   LoginSessionCookieName,
					Path:   "/",
					MaxAge: -1,
				})
				return next(c)
			}

			// Store user token in context for subsequent handlers
			c.Set(UserTokenKey, userToken.Token)
			c.Logger().Debug("User token stored in context")

			return next(c)
		}
	}
}

// GetUserToken retrieves the user token from context
func GetUserToken(c echo.Context) string {
	if token, ok := c.Get(UserTokenKey).(string); ok {
		return token
	}
	return ""
}

// IsAuthenticated checks if the current request is authenticated
func IsAuthenticated(c echo.Context) bool {
	return GetUserToken(c) != ""
}
