package client

import (
	"fmt"
	"time"

	config "github.com/oj-lab/reborn/configs"
	"github.com/oj-lab/user-service/pkg/userpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

// AuthServiceClient manages the gRPC client for auth service
type AuthServiceClient struct {
	client userpb.AuthServiceClient
	conn   *grpc.ClientConn
	config config.AuthServiceConfig
}

// NewAuthServiceClient creates a new auth service client with configuration
func NewAuthServiceClient(cfg config.AuthServiceConfig) (*AuthServiceClient, error) {
	opts := []grpc.DialOption{
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                30 * time.Second, // Send pings every 30 seconds instead of 10
			Timeout:             5 * time.Second,  // Wait 5 seconds for ping response
			PermitWithoutStream: false,            // Only send pings when there are active streams
		}),
	}

	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	// Attempt connection with retries
	var conn *grpc.ClientConn
	var err error

	conn, err = grpc.NewClient(cfg.Address, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to auth service at %s: %w", cfg.Address, err)
	}

	client := userpb.NewAuthServiceClient(conn)
	return &AuthServiceClient{
		client: client,
		conn:   conn,
		config: cfg,
	}, nil
}

// Close closes the gRPC connection
func (c *AuthServiceClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// GetClient returns the underlying gRPC client
func (c *AuthServiceClient) GetClient() userpb.AuthServiceClient {
	return c.client
}

// GetUserServiceClient returns a user service client using the same connection
func (c *AuthServiceClient) GetUserServiceClient() userpb.UserServiceClient {
	return userpb.NewUserServiceClient(c.conn)
}
