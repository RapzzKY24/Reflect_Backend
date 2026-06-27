package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort string

	DBHost     			string
	DBPort     			string
	DBUser     			string
	DBPassword 			string
	DBName     			string
	DBSSLMode  			string

	CloudinaryCloudName string
	CloudinaryAPIKey    string
	CloudinaryAPISecret string

	JWTSecret 			string

	RateLimitRegister  int
	RateLimitLogin     int
	RateLimitCheckout  int
	RateLimitUpload    int
	RateLimitAdmin     int
	RateLimitBurst     int
}

func LoadConfig() *Config {
	err := godotenv.Load()

	if err != nil {
		log.Println("No .env file found")
	}

	return &Config{
		AppPort: getEnv("APP_PORT", "8080"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "reflect_db"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		CloudinaryCloudName: getEnv("CLOUDINARY_CLOUD_NAME", ""),
		CloudinaryAPIKey:    getEnv("CLOUDINARY_API_KEY", ""),
		CloudinaryAPISecret: getEnv("CLOUDINARY_API_SECRET", ""),

		JWTSecret: getEnv("JWT_SECRET", "secret"),

		RateLimitRegister:  getEnvInt("RATE_LIMIT_REGISTER", 5),
		RateLimitLogin:     getEnvInt("RATE_LIMIT_LOGIN", 10),
		RateLimitCheckout:  getEnvInt("RATE_LIMIT_CHECKOUT", 5),
		RateLimitUpload:    getEnvInt("RATE_LIMIT_UPLOAD", 10),
		RateLimitAdmin:     getEnvInt("RATE_LIMIT_ADMIN", 30),
		RateLimitBurst:     getEnvInt("RATE_LIMIT_BURST", 2),
	}
}

func getEnv(key string, fallback string) string {
	value, exists := os.LookupEnv(key)

	if !exists {
		return fallback
	}

	return value
}

func getEnvInt(key string, fallback int) int {
	value, exists := os.LookupEnv(key)

	if !exists {
		return fallback
	}

	intVal, err := strconv.Atoi(value)

	if err != nil {
		return fallback
	}

	return intVal
}