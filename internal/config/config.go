package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost       string
	Port             string
	Mode             string
	GoogleConfigPath string
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()
	return Config{
		PublicHost:       getEnv("PUBLIC_HOST", "localhost"),
		Port:             getEnv("PORT", "8080"),
		Mode:             getEnv("MODE", "development"),
		GoogleConfigPath: getEnv("GOOGLE_CONFIG_PATH", "key.json"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		v, err := strconv.Atoi(value)
		if err != nil {
			log.Fatalf("environment variable %s could not be converted to type int \n error: %v", key, err)
		}
		return v
	}
	return fallback
}
