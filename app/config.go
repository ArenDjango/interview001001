//nolint:errcheck
package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	otellib "gitlab.b2broker.tech/pbsr/pbsr/backend/go/libs/otel"
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
	Trace    otellib.TraceConfig
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

	natsPort, err := strconv.Atoi(getEnv("NATS_PORT", "4222"))
	if err != nil {
		fmt.Println("Error converting  nats port to int type:", err)
	}

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
		ExternalServices: struct {
			Avanpost struct {
				Host                string
				Port                string
				OnboardingGroupUuid string
			}
		}{
			Avanpost: struct {
				Host                string
				Port                string
				OnboardingGroupUuid string
			}{
				Host:                getEnv("AVANPOST_SERVICE_HOST", ""),
				Port:                getEnv("AVANPOST_SERVICE_PORT", ""),
				OnboardingGroupUuid: getEnv("AVANPOST_ONBOARDING_GROUP_UUID", ""),
			},
		},
		Nats: struct {
			Host                  string
			Port                  int
			JwtCredentialFilePath string
		}{
			Host:                  getEnv("NATS_HOST", ""),
			Port:                  natsPort,
			JwtCredentialFilePath: getEnv("NATS_JWT_CREDENTIAL_FILE_PATH", ""),
		},
		Trace: otellib.TraceConfig{
			EnableTracing:     getEnv("ENABLE_TRACING", "false") == "true",
			CollectorEndpoint: getEnv("TRACE_COLLECTOR_ENDPOINT", ""),
			AppName:           getEnv("APP_NAME", ""),
		},
		MockGRPC: getEnv("MOCK_GRPC", "false"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
