package config

import (
	"errors"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	BotToken  string
	DB        *DB
	RedisDB   *RedisDB
	TextsPath string
}

type RedisDB struct {
	Host string
	Port string
}

type DB struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func LoadConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			log.Fatalf("config file not found: %v", err)
		}

		log.Fatalf("config file not read: %v", err)
	}

	var c Config
	err := viper.Unmarshal(&c)
	if err != nil {
		log.Fatalf("failed marshal config: %v", err)
	}

	return &c
}
