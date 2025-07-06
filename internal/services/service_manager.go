package services

import (
	"fmt"
	"log"
	"sync"

	config "github.com/oj-lab/reborn/configs"
)

// ServiceManager manages all application services
type ServiceManager struct {
	authService     *AuthService
	databaseService *DatabaseService
	mu              sync.RWMutex
}

// NewServiceManager creates a new service manager instance
func NewServiceManager() *ServiceManager {
	return &ServiceManager{}
}

// Initialize sets up all services with configuration
func (sm *ServiceManager) Initialize(cfg config.Config) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// Initialize database service first
	sm.databaseService = NewDatabaseService()
	if err := sm.databaseService.Initialize(cfg.Database); err != nil {
		return fmt.Errorf("failed to initialize database service: %w", err)
	}

	// Initialize auth service
	sm.authService = NewAuthService()
	if err := sm.authService.Initialize(cfg.AuthService); err != nil {
		return fmt.Errorf("failed to initialize auth service: %w", err)
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

// GetDatabaseService returns the database service instance
func (sm *ServiceManager) GetDatabaseService() *DatabaseService {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.databaseService
}

// Shutdown gracefully shuts down all services
func (sm *ServiceManager) Shutdown() error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	var lastErr error

	// Close database service
	if sm.databaseService != nil {
		if err := sm.databaseService.Close(); err != nil {
			log.Printf("Error closing database service: %v", err)
			lastErr = err
		}
	}

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

	// Check database service health
	if sm.databaseService != nil {
		health["database_service"] = sm.databaseService.IsHealthy()
	} else {
		health["database_service"] = false
	}

	// Check auth service health
	if sm.authService != nil {
		health["auth_service"] = sm.authService.IsHealthy()
	} else {
		health["auth_service"] = false
	}

	return health
}
