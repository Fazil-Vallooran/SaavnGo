package config

import (
	"os"
)

type Config struct {
	ServerPort      string
	JioSaavnBaseURL string
	DecryptionKey   string
}

func LoadConfig() *Config {
	return &Config{
		ServerPort:      getEnv("SERVER_PORT", "8080"),
		JioSaavnBaseURL: getEnv("JIOSAAVN_BASE_URL", "https://www.jiosaavn.com/api.php"),
		DecryptionKey:   getEnv("DECRYPTION_KEY", "38346591"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
