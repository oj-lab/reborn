package services

import (
	"testing"

	config "github.com/oj-lab/reborn/configs"
)

func TestDatabaseService_PostgreSQL(t *testing.T) {
	cfg := config.DatabaseConfig{
		Driver:   "postgres",
		Host:     "localhost",
		Port:     5432,
		Username: "user",
		Password: "password",
		Name:     "testdb",
		SSLMode:  "disable",
	}

	service := NewDatabaseService()
	
	// This should not create any directories and should not fail
	err := service.ensureSQLiteDirectoryExists(cfg)
	if err != nil {
		t.Errorf("PostgreSQL configuration should not cause errors: %v", err)
	}
}

func TestDatabaseService_NonSQLiteDrivers(t *testing.T) {
	drivers := []string{"postgres", "mysql", "POSTGRES", "MYSQL", "PostgreSQL"}
	
	for _, driver := range drivers {
		t.Run(driver, func(t *testing.T) {
			cfg := config.DatabaseConfig{
				Driver: driver,
				Name:   "test/path/database.db",
			}

			service := NewDatabaseService()
			err := service.ensureSQLiteDirectoryExists(cfg)
			
			if err != nil {
				t.Errorf("Non-SQLite driver %s should not cause errors: %v", driver, err)
			}
		})
	}
}