package app

import (
	"context"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

/*
All config should be required.
Optional only allowed if zero value of the type is expected being the default value.
time.Duration units are "ns", "us" (or "µs"), "ms", "s", "m", "h". as in time.ParseDuration().
*/

type (
	Postgres struct {
		ConnURI            string        `mapstructure:"POSTGRES_CONN_URI" validate:"required"`
		MaxPoolSize        int           `mapstructure:"POSTGRES_MAX_POOL_SIZE"`
		MaxIdleConnections int           `mapstructure:"POSTGRES_MAX_IDLE_CONNECTIONS"`
		MaxIdleTime        time.Duration `mapstructure:"POSTGRES_MAX_IDLE_TIME"`
		MaxLifeTime        time.Duration `mapstructure:"POSTGRES_MAX_LIFE_TIME"`
	}

	JWT struct {
		PrivateKeyPath string `mapstructure:"JWT_PRIVATE_KEY_PATH" validate:"required"`
		PublicKeyPath  string `mapstructure:"JWT_PUBLIC_KEY_PATH" validate:"required"`
	}

	Configuration struct {
		ServiceName string      `mapstructure:"SERVICE_NAME"`
		Postgres    Postgres    `mapstructure:",squash"`
		JWT         JWT         `mapstructure:",squash"`
		Translation Translation `mapstructure:",squash"`
		Environment string      `mapstructure:"ENV" validate:"required,oneof=development staging production"`
		BindAddress int         `mapstructure:"BIND_ADDRESS" validate:"required"`
		LogLevel    int         `mapstructure:"LOG_LEVEL" validate:"required"`
	}
)

func InitConfig(ctx context.Context) (*Configuration, error) {
	var cfg Configuration

	viper.SetConfigType("env")
	envFile := os.Getenv("ENV_FILE")
	if envFile == "" {
		envFile = ".env"
	}
	viper.SetConfigFile(envFile)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	viper.AutomaticEnv()

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
