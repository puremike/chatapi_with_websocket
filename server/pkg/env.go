package pkg

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func GetString(key, defaultValue string) string {
	if value, exist := os.LookupEnv(key); exist {
		return value
	}
	return defaultValue
}

func GetInt(key string, defaultValue int) int {
	if value, exist := os.LookupEnv(key); exist {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func GetDuration(key string, defaultValue time.Duration) time.Duration {
	if value, exist := os.LookupEnv(key); exist {
		if intTDuration, err := time.ParseDuration(value); err == nil {
			return intTDuration
		}
	}
	return defaultValue
}
