package config

import (
	"log"
	"os"
)

type Config struct {
	MQTTBroker string
	APIToken   string
	APIServer  string
	ClientID   string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	cfg := &Config{
		MQTTBroker: getEnv("MQTT_BROKER", "tcp://localhost:1883"),
		APIToken:   getEnv("API_TOKEN", ""),
		APIServer:  getEnv("API_SERVER", "https://example.com/api/temperature-data"),
		ClientID:   getEnv("CLIENT_ID", "raspberry-pi-hub"),
	}

	if cfg.APIToken == "" {
		log.Print("API_TOKEN is required")
	}

	return cfg
}

// getEnv gets an environment variable or a default value
func getEnv(key, defaultVal string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultVal
}
