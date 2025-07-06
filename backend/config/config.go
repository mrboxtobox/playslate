// Environment-based configuration following 12-factor principles
package config

import (
	"os"
	"strings"
)

type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
	
	// Frontend
	FrontendURL string
	CORSOrigins []string
	
	// Stripe
	StripeSecretKey   string
	StripeMonthlyPID  string
	StripeYearlyPID   string
	StripeWebhookSec  string
	
	// Email
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	
	// Google OAuth
	GoogleClientID string
}

func Load() *Config {
	return &Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		JWTSecret:   getEnv("JWT_SECRET", "dev-secret-change-in-production"),
		
		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:5173"),
		CORSOrigins: strings.Split(getEnv("CORS_ORIGINS", "http://localhost:5173,http://localhost:3000"), ","),
		
		StripeSecretKey:  getEnv("STRIPE_SECRET_KEY", ""),
		StripeMonthlyPID: getEnv("STRIPE_MONTHLY_PRICE_ID", ""),
		StripeYearlyPID:  getEnv("STRIPE_YEARLY_PRICE_ID", ""),
		StripeWebhookSec: getEnv("STRIPE_WEBHOOK_SECRET", ""),
		
		SMTPHost:     getEnv("SMTP_HOST", ""),
		SMTPPort:     getEnv("SMTP_PORT", "587"),
		SMTPUsername: getEnv("SMTP_USERNAME", ""),
		SMTPPassword: getEnv("SMTP_PASSWORD", ""),
		
		GoogleClientID: getEnv("GOOGLE_CLIENT_ID", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}