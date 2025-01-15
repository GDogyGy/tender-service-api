package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string `yaml:"ENV" env:"ENV" env-default:"local" env-required:"true"`
	HTTPServer `yaml:"HTTP_SERVER"`
	Postgres   `yaml:"POSTGRES"`
	DebugLevel string `yaml:"DEBUG_LEVEL" env:"DEBUG_LEVEL" env-default:"info"`
}

type HTTPServer struct {
	Address     string        `yaml:"SERVER_ADDRESS" env:"SERVER_ADDRESS" env-default:"localhost:8080" env-required:"true"`
	Timeout     time.Duration `yaml:"TIMEOUT" env:"TIMEOUT" env-default:"6s"`
	IdleTimeout time.Duration `yaml:"IDLE_TIMEOUT" env:"IDLE_TIMEOUT" env-default:"60s"`
}

type Postgres struct {
	PostgresConn     string `yaml:"POSTGRES_CONN" env:"POSTGRES_CONN" env-default:"postgres://root:123@localhost:5432/TenderApi" env-required:"true"`
	PostgresJdbcUrl  string `yaml:"POSTGRES_JDBC_URL" env:"POSTGRES_JDBC_URL" env-default:"jdbc:postgresql://root:123@localhost:5432/TenderApi" env-required:"true"`
	PostgresUsername string `yaml:"POSTGRES_USERNAME" env:"POSTGRES_USERNAME" env-default:"root" env-required:"true"`
	PostgresPassword string `yaml:"POSTGRES_PASSWORD" env:"POSTGRES_PASSWORD" env-default:"123" env-required:"true"`
	PostgresHost     string `yaml:"POSTGRES_HOST" env:"POSTGRES_HOST" env-default:"localhost" env-required:"true"`
	PostgresPort     string `yaml:"POSTGRES_PORT" env:"POSTGRES_PORT" env-default:"5432" env-required:"true"`
	PostgresDatabase string `yaml:"POSTGRES_DATABASE" env:"POSTGRES_DATABASE" env-default:"TenderApi" env-required:"true"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		log.Fatal("CONFIG_PATH id not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("CONFIG_PATH is not set: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Cannot read config: %s", err)
	}

	return &cfg
}
