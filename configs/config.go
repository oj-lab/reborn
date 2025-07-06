package config

import "github.com/oj-lab/go-webmods/app"

// Configuration keys constants
const (
	ServerPortKey         = "server.port"
	AuthServiceAddressKey = "auth_service.address"
	WebsiteDistPathKey    = "website.dist_path"
)

type Config struct {
	Server      ServerConfig
	AuthService AuthServiceConfig
	Website     WebsiteConfig
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
	}
	return cfg
}
