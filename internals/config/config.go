package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Port     string `mapstructure:"PORT"`
	Database string `mapstructure:"DATABASE_URL"`
	RedisURL string `mapstructure:"REDIS_URL"`
	AIAPIKey string `mapstructure:"AI_API_KEY"`
}

func LoadConfig() (Config, error) {
	var config Config

	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		log.Println("⚠️ No .env file found, falling back to environment variables")
	}

	if err := viper.Unmarshal(&config); err != nil {
		return config, err
	}

	return config, nil
}
