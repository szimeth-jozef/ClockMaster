package env

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type EnvConfig struct {
	Port   string  `envconfig:"PORT" required:"true"`
	DbPath *string `envconfig:"DB_PATH" required:"false"`
}

func New(filename string) (*EnvConfig, error) {
	err := godotenv.Load(filename)
	if err != nil {
		return nil, fmt.Errorf("Error loading .env file: %w", err)
	}

	var env EnvConfig
	err = envconfig.Process("", &env)
	if err != nil {
		return nil, fmt.Errorf("Error processing env config: %w", err)
	}

	return &env, nil
}
