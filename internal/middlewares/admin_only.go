package middlewares

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/oj-lab/reborn/internal/services"
	"github.com/oj-lab/user-service/pkg/userpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// AdminOnly returns a middleware that checks if the user is an admin
// This middleware should be used after LoginSession middleware
func AdminOnly(authService *services.AuthService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Check if user is authenticated
			if !IsAuthenticated(c) {
				return echo.NewHTTPError(http.StatusUnauthorized, "Authentication required")
			}

			// Get user token from context
			userToken := GetUserToken(c)
			if userToken == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid user token")
			}

			// Check if auth service is available
			if authService == nil || !authService.IsHealthy() {
				return echo.NewHTTPError(http.StatusServiceUnavailable, "User service unavailable")
			}

			// Get auth service client
			authClient := authService.GetClient()
			if authClient == nil {
				return echo.NewHTTPError(http.StatusServiceUnavailable, "User service client unavailable")
			}

			// Create user service client using the same connection
			userServiceClient := authClient.GetUserServiceClient()

			// Create context with user token for authentication
			md := metadata.Pairs("authorization", "Bearer "+userToken)
			ctx := metadata.NewOutgoingContext(c.Request().Context(), md)

			// Get current user information to check role
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

			// Check if user is admin
			if user.Role != userpb.UserRole_ADMIN {
				return echo.NewHTTPError(http.StatusForbidden, "Admin access required")
			}

			// Store user info in context for subsequent handlers
			c.Set("current_user", user)
			c.Logger().Debug("Admin user authenticated")

			return next(c)
		}
	}
}

// GetCurrentUser retrieves the current user from context
func GetCurrentUser(c echo.Context) *userpb.User {
	if user, ok := c.Get("current_user").(*userpb.User); ok {
		return user
	}
	return nil
}

// IsAdmin checks if the current user is an admin
func IsAdmin(c echo.Context) bool {
	user := GetCurrentUser(c)
	return user != nil && user.Role == userpb.UserRole_ADMIN
}
