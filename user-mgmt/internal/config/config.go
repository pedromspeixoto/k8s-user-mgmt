package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"go.uber.org/fx"
)

type Config struct {
	// Generic
	Environment  string `envconfig:"ENV" default:"development"`
	Port         string `envconfig:"APP_PORT" default:"8080"`
	AllowedHosts string `envconfig:"ALLOWED_HOSTS" default:"*"`

	// Logging
	LoggerType  string `envconfig:"LOGGER_TYPE" required:"false" default:"zap"`
	LoggerLevel int    `envconfig:"LOGGER_LEVEL" required:"false" default:"1"`

	// MySQL (Internal)
	MySQLHost     string `envconfig:"MYSQL_HOST" required:"false" default:"0.0.0.0"`
	MySQLPort     string `envconfig:"MYSQL_PORT" required:"false" default:"3306"`
	MySQLUser     string `envconfig:"MYSQL_USER" required:"false" default:"username"`
	MySQLPassword string `envconfig:"MYSQL_PASSWORD" required:"false" default:"password"`
	MySQLDBName   string `envconfig:"MYSQL_DB_NAME" required:"false" default:"dev_users"`

	// File Serving URL
	FileServingUrl string `envconfig:"FILE_SERVING_URL" required:"false" default:"http://localhost:3000"`
}

func ProvideConfig(cfgFile string) fx.Option {
	return fx.Provide(func() (*Config, error) {
		cfg, err := loadConfig(cfgFile)
		return cfg, err
	})
}

func loadConfig(cfgFile string) (cfg *Config, err error) {
	if cfgFile != "" {
		if cfg, err = readCfgFromFile(cfgFile); err != nil {
			return nil, err
		}
		return cfg, nil
	}
	if cfg, err = readCfgFromEnv(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func readCfgFromFile(cfgFile string) (*Config, error) {
	if err := godotenv.Load(cfgFile); err != nil {
		return nil, errors.WithStack(err)
	}
	return readCfgFromEnv()
}

func readCfgFromEnv() (*Config, error) {
	cfg := Config{}
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, errors.WithStack(err)
	}
	return &cfg, nil
}

func (c *Config) MySQLUrl() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		c.MySQLUser,
		c.MySQLPassword,
		c.MySQLHost,
		c.MySQLPort,
		c.MySQLDBName,
	)
}
