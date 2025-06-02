package gorm_client

import (
	"fmt"
	"sync"

	"github.com/oj-lab/reborn/common/app"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	configKeyDatabaseDriver = "gorm_client.database.driver"
	configKeyDatabaseHost   = "gorm_client.database.host"
	configKeyDatabasePort   = "gorm_client.database.port"
	configKeyDatabaseUser   = "gorm_client.database.username"
	configKeyDatabaseName   = "gorm_client.database.name"
	configKeyDatabasePass   = "gorm_client.database.password"
	configKeyDatabaseSSL    = "gorm_client.database.sslmode"
)

var (
	initMutx sync.Mutex
	db       *gorm.DB
)

func GetDB() *gorm.DB {
	if db == nil {
		initMutx.Lock()
		defer initMutx.Unlock()
		if db != nil {
			return db
		}

		driver := app.Config().GetString(configKeyDatabaseDriver)
		switch driver {
		case "postgres":
			db, err := openPostgres()
			if err != nil {
				panic(err)
			}
			return db
		default:
			panic(fmt.Sprintf("unsupported database driver: %s", driver))
		}
	}
	return db
}

func openPostgres() (db *gorm.DB, err error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		app.Config().GetString(configKeyDatabaseHost),
		app.Config().GetString(configKeyDatabasePort),
		app.Config().GetString(configKeyDatabaseUser),
		app.Config().GetString(configKeyDatabaseName),
		app.Config().GetString(configKeyDatabasePass),
		app.Config().GetString(configKeyDatabaseSSL),
	)
	db, err = gorm.Open(postgres.Open(dsn))
	if err != nil {
		return nil, err
	}
	return db, nil
}
