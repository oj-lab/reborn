package services

import (
	"os"
	"path/filepath"
	"testing"

	config "github.com/oj-lab/reborn/configs"
)

func TestEnsureSQLiteDirectoryExists(t *testing.T) {
	service := NewDatabaseService()

	tests := []struct {
		name          string
		cfg           config.DatabaseConfig
		expectError   bool
		shouldCreate  bool
		setupCleanup  func() (string, func())
	}{
		{
			name: "SQLite with nested directory path",
			cfg: config.DatabaseConfig{
				Driver: "sqlite",
				Name:   "test/data/test.db",
			},
			expectError:  false,
			shouldCreate: true,
			setupCleanup: func() (string, func()) {
				testDir := "test"
				return testDir, func() {
					os.RemoveAll(testDir)
				}
			},
		},
		{
			name: "SQLite with current directory",
			cfg: config.DatabaseConfig{
				Driver: "sqlite",
				Name:   "test.db",
			},
			expectError:  false,
			shouldCreate: false,
			setupCleanup: func() (string, func()) {
				return "", func() {}
			},
		},
		{
			name: "SQLite with existing directory",
			cfg: config.DatabaseConfig{
				Driver: "sqlite",
				Name:   "existing/test.db",
			},
			expectError:  false,
			shouldCreate: false,
			setupCleanup: func() (string, func()) {
				testDir := "existing"
				os.MkdirAll(testDir, 0755)
				return testDir, func() {
					os.RemoveAll(testDir)
				}
			},
		},
		{
			name: "PostgreSQL driver should be ignored",
			cfg: config.DatabaseConfig{
				Driver: "postgres",
				Name:   "test/data/test.db",
			},
			expectError:  false,
			shouldCreate: false,
			setupCleanup: func() (string, func()) {
				return "", func() {}
			},
		},
		{
			name: "Empty database name",
			cfg: config.DatabaseConfig{
				Driver: "sqlite",
				Name:   "",
			},
			expectError:  false,
			shouldCreate: false,
			setupCleanup: func() (string, func()) {
				return "", func() {}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDir, cleanup := tt.setupCleanup()
			defer cleanup()

			err := service.ensureSQLiteDirectoryExists(tt.cfg)

			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if tt.shouldCreate && tt.cfg.Name != "" {
				dir := filepath.Dir(tt.cfg.Name)
				if dir != "." {
					if _, err := os.Stat(dir); os.IsNotExist(err) {
						t.Errorf("Expected directory %s to be created but it doesn't exist", dir)
					}
				}
			}

			if testDir != "" && !tt.shouldCreate {
				// For PostgreSQL test, ensure no directory was created
				if tt.cfg.Driver != "sqlite" {
					if _, err := os.Stat(testDir); !os.IsNotExist(err) {
						t.Errorf("Directory %s should not have been created for non-SQLite driver", testDir)
					}
				}
			}
		})
	}
}

func TestDatabaseService_Initialize(t *testing.T) {
	tests := []struct {
		name         string
		cfg          config.DatabaseConfig
		expectError  bool
		setupCleanup func() func()
	}{
		{
			name: "SQLite initialization with directory creation",
			cfg: config.DatabaseConfig{
				Driver: "sqlite",
				Name:   "test_init/data/test.db",
			},
			expectError: false,
			setupCleanup: func() func() {
				return func() {
					os.RemoveAll("test_init")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := tt.setupCleanup()
			defer cleanup()

			service := NewDatabaseService()
			err := service.Initialize(tt.cfg)

			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			// Verify directory was created for SQLite
			if tt.cfg.Driver == "sqlite" && tt.cfg.Name != "" {
				dir := filepath.Dir(tt.cfg.Name)
				if dir != "." {
					if _, err := os.Stat(dir); os.IsNotExist(err) {
						t.Errorf("Expected directory %s to be created but it doesn't exist", dir)
					}
				}
			}

			// Clean up
			if service.db != nil {
				service.Close()
			}
		})
	}
}