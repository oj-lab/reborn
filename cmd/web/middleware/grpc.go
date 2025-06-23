package middleware

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/oj-lab/reborn/common/session"
	"google.golang.org/grpc/metadata"
)

func GetAuthenticatedContext(c echo.Context) (context.Context, error) {
	sess, ok := c.Get(ContextKeyUserSession).(*session.Session)
	if !ok || sess == nil {
		return nil, fmt.Errorf("invalid session from context")
	}

	if sess.JWT == "" {
		return nil, fmt.Errorf("jwt not found in session")
	}

	// Add token to gRPC context
	md := metadata.New(map[string]string{"authorization": "Bearer " + sess.JWT})
	ctx := metadata.NewOutgoingContext(c.Request().Context(), md)

	return ctx, nil
}
