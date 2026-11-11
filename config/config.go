package config

import (
	"os"
	"strings"
)

type Config struct {
	Port          string
	DBHost        string
	DBUser        string
	DBPass        string
	DBName        string
	DBPort        string
	JWTSecret     string
	AdminAPIKey   string
	AdminEmail    string
	AdminPassword string
	FrontendURL         string
	AutologinLinkSecret string
	ImportLeadsSecret   string
}

func Load() *Config {
	loadDotEnv()
	return &Config{
		Port:          getEnv("PORT", ":8080"),
		DBHost:        getEnv("DB_HOST", "127.0.0.1"),
		DBUser:        getEnv("DB_USER", "root"),
		DBPass:        getEnv("DB_PASS", "123456789"),
		DBName:        getEnv("DB_NAME", "godb"),
		DBPort:        getEnv("DB_PORT", "3306"),
		JWTSecret:     getEnv("JWT_SECRET", "change-me"),
		AdminAPIKey:   getEnv("ADMIN_API_KEY", ""),
		AdminEmail:    getEnv("ADMIN_EMAIL", ""),
		AdminPassword: getEnv("ADMIN_PASSWORD", ""),
		FrontendURL:         strings.TrimRight(getEnv("FRONTEND_URL", "https://autoshop.avinertech.com"), "/"),
		AutologinLinkSecret: getEnv("AUTOLOGIN_LINK_SECRET", ""),
		ImportLeadsSecret:   getEnv("IMPORT_LEADS_SECRET", ""),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
