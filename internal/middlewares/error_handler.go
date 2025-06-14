package middlewares

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// ErrorResponse represents a standard error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code"`
}

// ErrorHandler is a custom error handler for Echo
func ErrorHandler(err error, c echo.Context) {
	var (
		code = http.StatusInternalServerError
		msg  string
	)

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		if m, ok := he.Message.(string); ok {
			msg = m
		} else {
			msg = http.StatusText(code)
		}
	} else {
		msg = err.Error()
	}

	// Don't send the error message in production for security reasons
	if !c.Echo().Debug {
		switch code {
		case http.StatusNotFound:
			msg = "Resource not found"
		case http.StatusInternalServerError:
			msg = "Internal server error"
		case http.StatusUnauthorized:
			msg = "Unauthorized"
		case http.StatusForbidden:
			msg = "Forbidden"
		default:
			msg = http.StatusText(code)
		}
	}

	// Send response
	if !c.Response().Committed {
		if c.Request().Method == http.MethodHead {
			err = c.NoContent(code)
		} else {
			err = c.JSON(code, ErrorResponse{
				Error:   http.StatusText(code),
				Message: msg,
				Code:    code,
			})
		}
		if err != nil {
			c.Echo().Logger.Error(err)
		}
	}
}
