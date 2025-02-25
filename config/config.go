package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Env string

const (
	Env_Dev  Env = "dev"
	Env_Test Env = "test"
)

type Config struct {
	ApiServerHost           string `env:"APISERVER_HOST"`
	ApiServerPort           string `env:"APISERVER_PORT"`
	DBName                  string `env:"DB_NAME"`
	DBHost                  string `env:"DB_HOST"`
	DBPort                  string `env:"DB_PORT"`
	DBPortTest              string `env:"DB_PORT_TEST"`
	DBUser                  string `env:"DB_USER"`
	DBPassword              string `env:"DB_PASSWORD"`
	DBSSLMode               string `env:"DB_SSL_MODE"`
	DBSchema                string `env:"DB_SCHEMA"`
	Env                     Env    `env:"ENV" envDefault:"dev"`
	JwtSecret               string `env:"JWT_SECRET"`
	JwtAccessTokenLifetime  string `env:"JWT_ACCESS_TOKEN_LIFETIME"`
	JwtRefreshTokenLifetime string `env:"JWT_REFRESH_TOKEN_LIFETIME"`
	ProjectRoot             string `env:"PROJECT_ROOT"`
}

func New() (*Config, error) {
	cfg, err := env.ParseAs[Config]()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return &cfg, nil
}

func (c *Config) DataSourceName() string {

	if c.DBUser == "" {
		c.DBUser = "admin"
	}
	if c.DBPassword == "" {
		c.DBPassword = "secret"
	}
	if c.DBHost == "" {
		c.DBHost = "127.0.0.1"
	}
	if c.DBPort == "" {
		c.DBPort = "5432"
	}
	if c.DBName == "" {
		c.DBName = "asyncapi"
	}
	if c.DBSSLMode == "" {
		c.DBSSLMode = "disable"
	}
	if c.DBSchema == "" {
		c.DBSchema = "init_schema"
	}
	if c.Env == "" {
		c.Env = Env_Test
	}

	port := c.DBPort
	if c.Env == Env_Test {
		port = c.DBPortTest
	}

	// data sourse name
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.DBUser, c.DBPassword, c.DBHost, port, c.DBName, c.DBSSLMode)

	// dns := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s search_path=%s",
	// 	c.DBUser, c.DBPassword, c.DBHost, port, c.DBName, c.DBSSLMode, c.DBSchema)

	// fmt.Println("+++ [DataSourceName] dsn", dsn)
	return dsn
}
