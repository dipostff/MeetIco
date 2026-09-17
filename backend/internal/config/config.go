package config

import (
	"os"
)

type Config struct {
	DatabaseURL string
	JWTSecret   string
	SMTPHost    string
	SMTPPort    string
	SMTPUser    string
	SMTPPass    string
	SMTPFrom    string
	TelegramBotToken string
	BaseURL     string
	Port        string
}

func Load() *Config {
	return &Config{
		DatabaseURL:      getEnv("DATABASE_URL", ""),
		JWTSecret:        getEnv("JWT_SECRET", ""),
		SMTPHost:         getEnv("SMTP_HOST", ""),
		SMTPPort:         getEnv("SMTP_PORT", ""),
		SMTPUser:         getEnv("SMTP_USER", ""),
		SMTPPass:         getEnv("SMTP_PASS", ""),
		SMTPFrom:         getEnv("SMTP_FROM", ""),
		TelegramBotToken: getEnv("TELEGRAM_BOT_TOKEN", ""),
		BaseURL:          getEnv("BASE_URL", ""),
		Port:             getEnv("PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
