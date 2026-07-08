package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port           string
	DatabaseURL    string
	JWTSecret      string
	DefaultTimeout time.Duration
	SMTPHost       string
	SMTPPort       string
	SMTPUser       string
	SMTPPassword   string
	SMTPFrom       string
}

func Load() Config {
	return Config{
		Port:           env("PORT", "8080"),
		DatabaseURL:    env("DATABASE_URL", "postgres://peekaping:peekaping@localhost:5432/peekaping?sslmode=disable"),
		JWTSecret:      env("JWT_SECRET", "change-me-in-production"),
		DefaultTimeout: secondsEnv("DEFAULT_REQUEST_TIMEOUT_SECONDS", 10),
		SMTPHost:       env("SMTP_HOST", ""),
		SMTPPort:       env("SMTP_PORT", "587"),
		SMTPUser:       env("SMTP_USER", ""),
		SMTPPassword:   env("SMTP_PASSWORD", ""),
		SMTPFrom:       env("SMTP_FROM", "alerts@peekaping.local"),
	}
}

func env(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func secondsEnv(key string, fallback int) time.Duration {
	value, err := strconv.Atoi(env(key, ""))
	if err != nil || value <= 0 {
		value = fallback
	}
	return time.Duration(value) * time.Second
}
