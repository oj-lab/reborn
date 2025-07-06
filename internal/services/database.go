package services

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/oj-lab/go-webmods/gorm_client"
	config "github.com/oj-lab/reborn/configs"
	"gorm.io/gorm"
)

// DatabaseService manages database connections
type DatabaseService struct {
	db *gorm.DB
	mu sync.RWMutex
}

// NewDatabaseService creates a new DatabaseService instance
func NewDatabaseService() *DatabaseService {
	return &DatabaseService{}
}

// Initialize sets up the database connection with provided config
func (s *DatabaseService) Initialize(cfg config.DatabaseConfig) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// For SQLite, ensure directory exists
	if err := s.ensureSQLiteDirectoryExists(cfg); err != nil {
		return fmt.Errorf("failed to ensure SQLite directory exists: %w", err)
	}

	// Create gorm client config
	gormConfig := gorm_client.Config{
		Driver:   cfg.Driver,
		Host:     cfg.Host,
		Port:     cfg.Port,
		Username: cfg.Username,
		Password: cfg.Password,
		Name:     cfg.Name,
		SSLMode:  cfg.SSLMode,
	}

	// Initialize database connection
	db := gorm_client.NewDB(gormConfig)
	if db == nil {
		return fmt.Errorf("failed to initialize database connection")
	}

	s.db = db

	log.Printf("Database service initialized successfully with driver: %s", cfg.Driver)
	return nil
}

// ensureSQLiteDirectoryExists creates the directory for SQLite database file if it doesn't exist
func (s *DatabaseService) ensureSQLiteDirectoryExists(cfg config.DatabaseConfig) error {
	// Only handle SQLite driver
	if strings.ToLower(cfg.Driver) != "sqlite" {
		return nil
	}

	// Skip if Name is empty
	if cfg.Name == "" {
		return nil
	}

	// Get directory path from the database file path
	dbDir := filepath.Dir(cfg.Name)
	
	// Skip if it's current directory
	if dbDir == "." {
		return nil
	}

	// Check if directory exists
	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		// Create directory with appropriate permissions
		if err := os.MkdirAll(dbDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dbDir, err)
		}
		log.Printf("Created directory for SQLite database: %s", dbDir)
	}

	return nil
}

// GetDB returns the database connection
func (s *DatabaseService) GetDB() *gorm.DB {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.db
}

// Close closes the database connection
func (s *DatabaseService) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.db != nil {
		sqlDB, err := s.db.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

// IsHealthy checks if the database connection is available
func (s *DatabaseService) IsHealthy() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	if s.db == nil {
		return false
	}

	sqlDB, err := s.db.DB()
	if err != nil {
		return false
	}

	return sqlDB.Ping() == nil
}