package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBName        string
	DBAddress     string
	DBPassword    string
	DBUser        string
	AccessSecret  string
	RefreshSecret string
	Port          string
	AccessExpiry  int
	RefreshExpiry int
}

func LoadConfig() (Config, error) {
	err := godotenv.Load()
	cfg := Config{
		DBName:        os.Getenv("DB_NAME"),
		DBAddress:     os.Getenv("DB_ADDRESS"),
		DBPassword:    os.Getenv("DB_PASSWORD"),
		DBUser:        os.Getenv("DB_USER"),
		AccessSecret:  os.Getenv("ACCESS_SECRET"),
		RefreshSecret: os.Getenv("REFRESH_SECRET"),
		Port:          os.Getenv("PORT"),
		AccessExpiry:  3600,
		RefreshExpiry: 604800,
	}
	if err != nil {
		return cfg, err
	}
	if cfg.DBName == "" || cfg.DBPassword == "" || cfg.DBUser == "" {
		return cfg, fmt.Errorf("missing db environment variables")
	}
	if cfg.AccessSecret == "" || cfg.RefreshSecret == "" {
		return cfg, fmt.Errorf("missing jwt secret environment variables")
	}
	if cfg.DBAddress == "" {
		cfg.DBAddress = "localhost:3306"
	}
	if cfg.Port == "" {
		cfg.Port = "8080"
	}
	return cfg, nil
}
