package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/oj-lab/reborn/common/session"
)

var (
	sessionManager = session.NewManager()
)

func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("session_id")
		if err != nil {
			return c.String(http.StatusUnauthorized, "Missing session cookie")
		}
		sessionID := cookie.Value

		session, err := sessionManager.Get(c.Request().Context(), sessionID)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to get session")
		}

		if session == nil {
			return c.String(http.StatusUnauthorized, "Invalid session")
		}

		c.Set("userID", uint64(session.UserID))
		return next(c)
	}
}
