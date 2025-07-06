package config

import "github.com/oj-lab/go-webmods/app"

// Configuration keys constants
const (
	ServerPortKey         = "server.port"
	AuthServiceAddressKey = "auth_service.address"
	WebsiteDistPathKey    = "website.dist_path"
	DatabaseDriverKey     = "database.driver"
	DatabaseHostKey       = "database.host"
	DatabasePortKey       = "database.port"
	DatabaseUsernameKey   = "database.username"
	DatabasePasswordKey   = "database.password"
	DatabaseNameKey       = "database.name"
	DatabaseSSLModeKey    = "database.ssl_mode"
)

type Config struct {
	Server      ServerConfig
	AuthService AuthServiceConfig
	Website     WebsiteConfig
	Database    DatabaseConfig
}

type ServerConfig struct {
	Port uint
}

type AuthServiceConfig struct {
	Address string
}

type WebsiteConfig struct {
	DistPath string
}

type DatabaseConfig struct {
	Driver   string
	Host     string
	Port     int
	Username string
	Password string
	Name     string
	SSLMode  string
}

func Load() Config {
	cfg := Config{
		Server: ServerConfig{
			Port: app.Config().GetUint(ServerPortKey),
		},
		AuthService: AuthServiceConfig{
			Address: app.Config().GetString(AuthServiceAddressKey),
		},
		Website: WebsiteConfig{
			DistPath: app.Config().GetString(WebsiteDistPathKey),
		},
		Database: DatabaseConfig{
			Driver:   app.Config().GetString(DatabaseDriverKey),
			Host:     app.Config().GetString(DatabaseHostKey),
			Port:     app.Config().GetInt(DatabasePortKey),
			Username: app.Config().GetString(DatabaseUsernameKey),
			Password: app.Config().GetString(DatabasePasswordKey),
			Name:     app.Config().GetString(DatabaseNameKey),
			SSLMode:  app.Config().GetString(DatabaseSSLModeKey),
		},
	}
	return cfg
}
