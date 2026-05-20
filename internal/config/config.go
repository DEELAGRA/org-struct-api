package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	ServerPort int
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
}

func getEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatal("An empty string was passed")
	}
	return val
}

func getEnvINT(key string) (int, error) {
	val := os.Getenv(key)
	if val == "" {
		log.Fatal("An empty string was passed")
	}
	return strconv.Atoi(val)
}
func Load() (*Config, error) {
	port, err := getEnvINT("SERVER_PORT")
	if err != nil {
		return nil, err
	}

	dbPort, err := getEnvINT("DB_PORT")
	if err != nil {
		return nil, err
	}

	return &Config{
		ServerPort: port,
		DBHost:     getEnv("DB_HOST"),
		DBPort:     dbPort,
		DBUser:     getEnv("DB_USER"),
		DBPassword: getEnv("DB_PASSWORD"),
		DBName:     getEnv("DB_NAME"),
		DBSSLMode:  getEnv("DB_SSLMODE"),
	}, nil
}
