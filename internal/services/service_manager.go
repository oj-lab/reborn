package services

import (
	"log"
	"sync"

	config "github.com/oj-lab/reborn/configs"
)

// ServiceManager manages all application services
type ServiceManager struct {
	authService *AuthService
	mu          sync.RWMutex
}

// NewServiceManager creates a new service manager instance
func NewServiceManager() *ServiceManager {
	return &ServiceManager{}
}

// Initialize sets up all services with configuration
func (sm *ServiceManager) Initialize(cfg config.Config) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// Initialize auth service
	sm.authService = NewAuthService()
	if err := sm.authService.Initialize(cfg.AuthService); err != nil {
		return err
	}

	log.Println("Service manager initialized successfully")
	return nil
}

// GetAuthService returns the auth service instance
func (sm *ServiceManager) GetAuthService() *AuthService {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.authService
}

// Shutdown gracefully shuts down all services
func (sm *ServiceManager) Shutdown() error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	var lastErr error

	// Close auth service
	if sm.authService != nil {
		if err := sm.authService.Close(); err != nil {
			log.Printf("Error closing auth service: %v", err)
			lastErr = err
		}
	}

	log.Println("Service manager shutdown completed")
	return lastErr
}

// HealthCheck returns the health status of all services
func (sm *ServiceManager) HealthCheck() map[string]bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	health := make(map[string]bool)

	// Check auth service health
	if sm.authService != nil {
		health["auth_service"] = sm.authService.IsHealthy()
	} else {
		health["auth_service"] = false
	}

	return health
}
