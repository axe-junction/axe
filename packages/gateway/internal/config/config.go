package config

import (
	"os"
	"strconv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	OAuth    OAuthConfig
	Services ServicesConfig
}

type ServerConfig struct {
	Port string
	Host string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type OAuthConfig struct {
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
	JWTSecret          string
	SessionSecret      string
}

type ServicesConfig struct {
	RoutingService RoutingServiceConfig
	VTCService     VTCServiceConfig
}

type RoutingServiceConfig struct {
	Host string
	Port string
}

type VTCServiceConfig struct {
	Host string
	Port string
}

func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Host: getEnv("HOST", "0.0.0.0"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "gateway_db"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		OAuth: OAuthConfig{
			GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
			GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
			GoogleRedirectURL:  getEnv("GOOGLE_REDIRECT_URL", "http://localhost:8080/auth/google/callback"),
			JWTSecret:          getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-in-production"),
			SessionSecret:      getEnv("SESSION_SECRET", "your-super-secret-session-key-change-in-production"),
		},
		Services: ServicesConfig{
			RoutingService: RoutingServiceConfig{
				Host: getEnv("ROUTING_SERVICE_HOST", "localhost"),
				Port: getEnv("ROUTING_SERVICE_PORT", "50051"),
			},
			VTCService: VTCServiceConfig{
				Host: getEnv("VTC_SERVICE_HOST", "localhost"),
				Port: getEnv("VTC_SERVICE_PORT", "50052"),
			},
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
