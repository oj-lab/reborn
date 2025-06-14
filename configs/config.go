package config

import "github.com/oj-lab/go-webmods/app"

// Configuration keys constants
const (
	ServerPortKey         = "server.port"
	AuthServiceAddressKey = "auth_service.address"
)

type Config struct {
	Server      ServerConfig
	AuthService AuthServiceConfig
}

type ServerConfig struct {
	Port uint
}

type AuthServiceConfig struct {
	Address string
}

func Load() Config {
	cfg := Config{
		Server: ServerConfig{
			Port: app.Config().GetUint(ServerPortKey),
		},
		AuthService: AuthServiceConfig{
			Address: app.Config().GetString(AuthServiceAddressKey),
		},
	}
	return cfg
}
