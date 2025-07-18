package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server ServerConfig `json:"server"`
	Rates  RatesConfig  `json:"rates"`
}

type ServerConfig struct {
	Port           string        `json:"port"`
	ReadTimeout    time.Duration `json:"readTimeout"`
	WriteTimeout   time.Duration `json:"writeTimeout"`
	RequestTimeout time.Duration `json:"requestTimeout"`
}

type RatesConfig struct {
	ServiceA    float64 `json:"serviceA"`
	ServiceBMin float64 `json:"serviceBMin"`
	ServiceBMax float64 `json:"serviceBMax"`
	ServiceC    float64 `json:"serviceC"`
	ServiceD    float64 `json:"serviceD"`
}

func NewConfig() *Config {
	godotenv.Load()

	return &Config{
		Server: ServerConfig{
			Port:           getEnvOrDefault("PORT", ":8080"),
			ReadTimeout:    parseDurationOrDefault("READ_TIMEOUT", 15*time.Second),
			WriteTimeout:   parseDurationOrDefault("WRITE_TIMEOUT", 15*time.Second),
			RequestTimeout: parseDurationOrDefault("REQUEST_TIMEOUT", 30*time.Second),
		},
		Rates: RatesConfig{
			ServiceA:    parseFloatOrDefault("SERVICE_A_RATE", 60.0),
			ServiceBMin: parseFloatOrDefault("SERVICE_B_MIN_RATE", 50.0),
			ServiceBMax: parseFloatOrDefault("SERVICE_B_MAX_RATE", 70.0),
			ServiceC:    parseFloatOrDefault("SERVICE_C_RATE", 58.0),
			ServiceD:    parseFloatOrDefault("SERVICE_D_RATE", 62.0),
		},
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func parseFloatOrDefault(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseFloat(value, 64); err == nil {
			return parsed
		}
	}
	return defaultValue
}

func parseDurationOrDefault(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if parsed, err := time.ParseDuration(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}
