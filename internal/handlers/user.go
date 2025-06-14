package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/oj-lab/reborn/internal/middlewares"
	"github.com/oj-lab/reborn/internal/services"
	"github.com/oj-lab/user-service/pkg/userpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	authService *services.AuthService
}

// NewUserHandler creates a new user handler instance
func NewUserHandler(authService *services.AuthService) *UserHandler {
	return &UserHandler{
		authService: authService,
	}
}

// GetCurrentUser returns the current authenticated user information
//
//	@Summary		Get current user
//	@Description	Retrieve the information of the currently authenticated user
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	userpb.User
//	@Failure		401	{object}	echo.HTTPError	"Unauthorized"
//	@Failure		500	{object}	echo.HTTPError	"Internal Server Error"
//	@Router			/api/v1/user/me [get]
func (h *UserHandler) GetCurrentUser(c echo.Context) error {
	// Check if user is authenticated
	if !middlewares.IsAuthenticated(c) {
		return echo.NewHTTPError(http.StatusUnauthorized, "Authentication required")
	}

	// Get user token from context
	userToken := middlewares.GetUserToken(c)
	if userToken == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid user token")
	}

	// Check if auth service is available
	if h.authService == nil || !h.authService.IsHealthy() {
		return echo.NewHTTPError(http.StatusServiceUnavailable, "User service unavailable")
	}

	// Get auth service client
	authClient := h.authService.GetClient()
	if authClient == nil {
		return echo.NewHTTPError(http.StatusServiceUnavailable, "User service client unavailable")
	}

	// Create user service client using the same connection
	var userServiceClient userpb.UserServiceClient = authClient.GetUserServiceClient()
	// Create context with user token for authentication
	md := metadata.Pairs("authorization", "Bearer "+userToken)
	ctx := metadata.NewOutgoingContext(c.Request().Context(), md)
	// Get current user information
	user, err := userServiceClient.GetCurrentUser(ctx, &emptypb.Empty{})
	if err != nil {
		// Handle gRPC errors
		if grpcStatus, ok := status.FromError(err); ok {
			switch grpcStatus.Code() {
			case codes.Unauthenticated:
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid or expired token")
			case codes.NotFound:
				return echo.NewHTTPError(http.StatusNotFound, "User not found")
			case codes.PermissionDenied:
				return echo.NewHTTPError(http.StatusForbidden, "Permission denied")
			default:
				return echo.NewHTTPError(
					http.StatusInternalServerError,
					"Failed to get user information",
				)
			}
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get user information")
	}

	return c.JSON(http.StatusOK, user)
}
