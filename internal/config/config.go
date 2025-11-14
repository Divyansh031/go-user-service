package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env      string        `yaml:"env" env:"ENV" env-default:"development"`
	GRPC     GRPCConfig    `yaml:"grpc"`
	HTTP     HTTPConfig    `yaml:"http"`
	ScyllaDB ScyllaDBConfig `yaml:"scylladb"`
	Log      LogConfig     `yaml:"log"`
}

type GRPCConfig struct {
	Port int `yaml:"port" env:"GRPC_PORT" env-default:"50051"`
}

type HTTPConfig struct {
	Port int `yaml:"port" env:"HTTP_PORT" env-default:"8080"`
}

type ScyllaDBConfig struct {
	Hosts       []string `yaml:"hosts" env:"SCYLLA_HOSTS" env-default:"localhost"`
	Port        int      `yaml:"port" env:"SCYLLA_PORT" env-default:"9042"`
	Keyspace    string   `yaml:"keyspace" env:"SCYLLA_KEYSPACE" env-default:"userservice"`
	Consistency string   `yaml:"consistency" env:"SCYLLA_CONSISTENCY" env-default:"QUORUM"`
}

type LogConfig struct {
	Level string `yaml:"level" env:"LOG_LEVEL" env-default:"info"`
}

// Load loads configuration from file or environment
func Load() (*Config, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./config/config.yaml"
	}

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Use environment variables only
		return loadFromEnv()
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	return &cfg, nil
}

func loadFromEnv() (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("failed to read env: %w", err)
	}
	return &cfg, nil
}

// MustLoad loads config and panics on error
func MustLoad() *Config {
	cfg, err := Load()
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}
	return cfg
}