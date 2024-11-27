package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
	ServerAddress    string
	ServerPort       string
	LogLevel         string
	ExternalAPI      string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile("./.env")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	config := &Config{
		PostgresHost:     viper.GetString("POSTGRES_HOST"),
		PostgresPort:     viper.GetString("POSTGRES_PORT"),
		PostgresUser:     viper.GetString("POSTGRES_USER"),
		PostgresPassword: viper.GetString("POSTGRES_PASSWORD"),
		PostgresDB:       viper.GetString("POSTGRES_DB"),
		ServerAddress:    viper.GetString("SERVER_ADDRESS"),
		ServerPort:       viper.GetString("SERVER_PORT"),
		LogLevel:         viper.GetString("LOG_LEVEL"),
		ExternalAPI:      viper.GetString("EXTERNAL_API"),
	}
	return config, nil
}

func (cfg *Config) PostgresURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresDB)
}
