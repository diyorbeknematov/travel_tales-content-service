package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"
)

type Config struct {
	GRPC_PORT        string
	USER_CLIENT_PORT string
	DB_HOST          string
	DB_PORT          string
	DB_USER          string
	DB_NAME          string
	DB_PASSWORD      string
}

func Load() Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	config := Config{}

	config.DB_HOST = cast.ToString(coalesce("DB_HOST", "localhost"))
	config.DB_PORT = cast.ToString(coalesce("DB_PORT", "5432"))
	config.DB_USER = cast.ToString(coalesce("DB_USER", "postgres"))
	config.DB_NAME = cast.ToString(coalesce("DB_NAME", "postgres"))
	config.DB_PASSWORD = cast.ToString(coalesce("DB_PASSWORD", "passwrod"))
	config.USER_CLIENT_PORT = cast.ToString(coalesce("USER_CLIENT_PORT", 50050))
	config.GRPC_PORT = cast.ToString(coalesce("GRPC_PORT", 50051))

	return config
}

func coalesce(env string, defaultValue interface{}) interface{} {
	value, exists := os.LookupEnv(env)
	if !exists {
		return defaultValue
	}
	return value
}
