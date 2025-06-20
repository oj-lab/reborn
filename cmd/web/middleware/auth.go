package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/oj-lab/reborn/common/session"
)

const (
	CookieUserSessionID   = "user_session_id"
	ContextKeyUserSession = "user_session"
)

func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sessionManager := session.NewManager()
		cookie, err := c.Cookie(CookieUserSessionID)
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

		c.Set(ContextKeyUserSession, session)
		return next(c)
	}
}
