package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/oj-lab/reborn/internal/middlewares"
	"github.com/oj-lab/reborn/internal/services"
	"github.com/oj-lab/user-service/pkg/userpb"
)

// AuthHandler handles authentication related HTTP requests
type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler creates a new auth handler with injected auth service
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Login handles OAuth login requests
func (h *AuthHandler) Login(ctx echo.Context) error {
	// Check if auth service is available
	if !h.authService.IsHealthy() {
		return echo.NewHTTPError(http.StatusServiceUnavailable, "Auth service unavailable")
	}

	// Get provider from query parameter, default to github
	provider := ctx.QueryParam("provider")
	if provider == "" {
		provider = "github"
	}

	// Automatically generate redirect URL based on request scheme and host
	// E.g. When locally running, it will be http://localhost:8080/auth/callback
	// When deployed, it will be https://example.com/auth/callback
	redirectURL := fmt.Sprintf("%s://%s/auth/callback", ctx.Scheme(), ctx.Request().Host)

	// Get OAuth URL from auth service
	client := h.authService.GetClient()
	resp, err := client.GetClient().
		GetOAuthCodeURL(ctx.Request().Context(), &userpb.GetOAuthCodeURLRequest{
			Provider:    provider,
			RedirectUrl: &redirectURL,
		})
	if err != nil {
		slog.Error("Failed to get OAuth URL", "error", err, "provider", provider, "redirect_url", redirectURL)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get OAuth URL")
	}

	return ctx.Redirect(http.StatusFound, resp.GetUrl())
}

// Callback handles OAuth callback requests
func (h *AuthHandler) Callback(ctx echo.Context) error {
	code := ctx.QueryParam("code")
	if code == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing code parameter")
	}
	state := ctx.QueryParam("state")
	if state == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing state parameter")
	}

	client := h.authService.GetClient()
	resp, err := client.GetClient().
		LoginByOAuth(ctx.Request().Context(), &userpb.LoginByOAuthRequest{
			Code:  code,
			State: state,
		})
	if err != nil {
		slog.Error("Failed to login by OAuth", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to login by OAuth")
	}
	ctx.SetCookie(&http.Cookie{
		Name:   middlewares.LoginSessionCookieName,
		Value:  resp.GetId(),
		Path:   "/",
		MaxAge: int(time.Until(resp.ExpiresAt.AsTime()).Seconds()),
	})
	return ctx.Redirect(http.StatusFound, "/") // Redirect to home page after login
}

// Logout handles logout requests
func (h *AuthHandler) Logout(ctx echo.Context) error {
	// Clear the login session cookie
	ctx.SetCookie(&http.Cookie{
		Name:   middlewares.LoginSessionCookieName,
		Value:  "",
		Path:   "/",
		MaxAge: -1, // Expire immediately
	})
	return ctx.Redirect(http.StatusFound, "/") // Redirect to home page after logout
}
