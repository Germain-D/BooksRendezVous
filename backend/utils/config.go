package utils

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

// Config holds all the environment variables.
type Config struct {

	// Database
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string

	// Frontend
	FrontendURL string

	// Server
	ServerPort string

	// Secret Key
	SecretKey string

	//log level
	LogLevel string

	SMTPHost     string `mapstructure:"SMTP_HOST"`
	SMTPPort     string `mapstructure:"SMTP_PORT"`
	SMTPUsername string `mapstructure:"SMTP_USERNAME"`
	SMTPPassword string `mapstructure:"SMTP_PASSWORD"`
	FromEmail    string `mapstructure:"FROM_EMAIL"`

	// Authorized emails for registration
	AuthorizedEmails []string
}

// Initialize a global SugaredLogger
var SugaredLogger *zap.SugaredLogger

func init() {
	// Initialize the logger
	logger, _ := zap.NewProduction()
	SugaredLogger = logger.Sugar()
}

func LoadConfig() (*Config, error) {
	sugar := SugaredLogger

	if err := godotenv.Load("./.env"); err != nil {
		sugar.Warnw("Warning: .env file not found",
			"error", err,
		)
	}

	return &Config{

		// Database
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "bibliotheque"),
		DBPort:     getEnv("DB_PORT", "5432"),

		// Frontend
		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:3000"),

		// Server
		ServerPort: getEnv("SERVER_PORT", ":8000"),

		// Secret Key
		SecretKey: getEnv("SECRET_KEY", "secret"),

		//log level
		LogLevel: getEnv("LOG_LEVEL", "info"),

		SMTPHost:     getEnv("SMTP_HOST", "ssl0.ovh.net"),
		SMTPPort:     getEnv("SMTP_PORT", "2525"),
		SMTPUsername: getEnv("SMTP_USER", "support@booksrendezvous.fr"),
		SMTPPassword: getEnv("SMTP_PASSWORD", ""),
		FromEmail:    getEnv("FROM_EMAIL", "support@booksrendezvous.fr"),

		// Authorized emails for registration (comma-separated)
		AuthorizedEmails: getEnvAsStringSlice("AUTHORIZED_EMAILS", []string{"john@example.com"}),
	}, nil
}

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		SugaredLogger.Debugw("Using default value for config",
			"key", key,
			"default", defaultValue,
		)
		return defaultValue
	}
	return value
}

// getEnvAsInt retrieves an environment variable as an integer or returns a default value.
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("Invalid value for %s, using default: %d\n", key, defaultValue)
		return defaultValue
	}
	return value
}

// getEnvAsStringSlice retrieves an environment variable as a slice of strings (comma-separated) or returns a default value.
func getEnvAsStringSlice(key string, defaultValue []string) []string {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}

	// Split by comma and trim whitespace
	parts := strings.Split(valueStr, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	if len(result) == 0 {
		return defaultValue
	}

	return result
}
