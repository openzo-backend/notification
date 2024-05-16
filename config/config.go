package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	MODE             string `mapstructure:"MODE"`
	DB_URL           string `mapstructure:"DB_URL"`
	HTTPPort         string `mapstructure:"HTTP_PORT"`
	GRPCPort         string `mapstructure:"GRPC_PORT"`
	UserGrpc         string `mapstructure:"USER_GRPC"`
	ImageGrpc        string `mapstructure:"IMAGE_GRPC"`
	StoreGrpc        string `mapstructure:"STORE_GRPC"`
	NotificationGrpc string `mapstructure:"NOTIFICATION_GRPC"`
	ProductGrpc      string `mapstructure:"PRODUCT_GRPC"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile("./config/config.yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var Config Config
	if err := viper.Unmarshal(&Config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	return &Config, nil
}
