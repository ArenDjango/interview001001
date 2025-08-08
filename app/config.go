//nolint:errcheck
package app

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type PathToEnv struct{}

// Config represents the application configuration.
type Config struct {
	Env   string
	Debug bool
	Url   string
	DB    struct {
		Host      string
		Port      string
		User      string
		Password  string
		Database  string
		BatchSize int
	}
	ExternalServices struct {
		Avanpost struct {
			Host                string
			Port                string
			OnboardingGroupUuid string
		}
	}
	Nats struct {
		Host                  string
		Port                  int
		JwtCredentialFilePath string
	}
	MockGRPC string
}

//nolint:funlen
func LoadConfig(ctx context.Context) *Config {
	pathToEnv := ctx.Value(PathToEnv{})

	var err error

	if pathToEnv != nil {
		err = godotenv.Load(pathToEnv.(string))
	} else {
		err = godotenv.Load()
	}

	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	debug, _ := strconv.ParseBool(getEnv("DEBUG", "false"))
	batchSize, _ := strconv.Atoi(getEnv("DB_BATCH_SIZE", "100"))

	return &Config{
		Env:   getEnv("APP_ENV", "dev"),
		Debug: debug,
		Url:   getEnv("APP_URL", ""),
		DB: struct {
			Host      string
			Port      string
			User      string
			Password  string
			Database  string
			BatchSize int
		}{
			Host:      getEnv("DB_HOST", ""),
			Port:      getEnv("DB_PORT", ""),
			User:      getEnv("DB_USERNAME", ""),
			Password:  getEnv("DB_PASSWORD", ""),
			Database:  getEnv("DB_DATABASE", ""),
			BatchSize: batchSize,
		},
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
