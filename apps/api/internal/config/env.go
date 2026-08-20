package config

import "fmt"

type Environment string

const (
	EnvironmentDevelopment Environment = "development"
	EnvironmentProduction  Environment = "production"
)

func (e *Environment) UnmarshalText(text []byte) error {
	value := Environment(text)

	switch value {
	case EnvironmentDevelopment, EnvironmentProduction:
		*e = value
		return nil
	default:
		return fmt.Errorf("Invalid environment: %q", value)
	}
}

type Config struct {
	Environment Environment `env:"ENV,required"`
	DatabaseURL string      `env:"DATABASE_URL,required"`
}
