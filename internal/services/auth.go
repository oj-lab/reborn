package services

import (
	"log"
	"sync"

	config "github.com/oj-lab/reborn/configs"
	"github.com/oj-lab/reborn/internal/client"
)

// AuthService manages auth service client connections
type AuthService struct {
	client *client.AuthServiceClient
	mu     sync.RWMutex
}

// NewAuthService creates a new AuthService instance
func NewAuthService() *AuthService {
	return &AuthService{}
}

// Initialize sets up the auth service client with provided config
func (s *AuthService) Initialize(cfg config.AuthServiceConfig) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	client, err := client.NewAuthServiceClient(cfg)
	if err != nil {
		return err
	}

	s.client = client
	return nil
}

// GetClient returns the auth service client
func (s *AuthService) GetClient() *client.AuthServiceClient {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.client
}

// SetClient sets the auth service client (for testing purposes)
func (s *AuthService) SetClient(client *client.AuthServiceClient) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Close existing client if any
	if s.client != nil {
		if err := s.client.Close(); err != nil {
			log.Printf("Error closing existing auth service connection: %v", err)
		}
	}

	s.client = client
}

// Reconnect attempts to reconnect to the auth service
func (s *AuthService) Reconnect(cfg *config.AuthServiceConfig) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Close existing connection
	if s.client != nil {
		if err := s.client.Close(); err != nil {
			log.Printf("Error closing auth service connection during reconnect: %v", err)
		}
	}

	// Create new connection
	client, err := client.NewAuthServiceClient(*cfg)
	if err != nil {
		return err
	}

	s.client = client
	return nil
}

// Close closes the auth service connections
func (s *AuthService) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.client != nil {
		return s.client.Close()
	}
	return nil
}

// IsHealthy checks if the auth service client is available
func (s *AuthService) IsHealthy() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.client != nil
}
