package configs

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

// setup struct config
type Config struct {
	App AppConfig
	DB  DBConfig
}

// config untuk application
type AppConfig struct {
	Port        string
	Environment string
}

// config untuk database
type DBConfig struct {
	Host string
	Port string
	User string
	Pass string
	Name string
}

// config jwt
type TokenConfig struct {
	JWTSecret   []byte
	TokenExpiry time.Duration
}

var Token = TokenConfig{
	JWTSecret:   []byte("400225a3c1f2b8f6db67195ea664d805b3f4d5291ca3bd4294fed07dd0dc4f5d"),
	TokenExpiry: time.Hour * 24,
}

func LoadConfig() Config {
	return Config{
		App: AppConfig{
			Port:        getConfig("APP_PORT"),
			Environment: getConfig("APP_ENV"),
		},
		DB: DBConfig{
			Host: getConfig("DB_HOST"),
			Port: getConfig("DB_PORT"),
			User: getConfig("DB_USER"),
			Pass: getConfig("DB_PASS"),
			Name: getConfig("DB_NAME"),
		},
	}
}

func getConfig(key string) string {
	err := godotenv.Overload(".env")
	if err != nil {
		log.Println("error :", err.Error())
		return ""
	}
	return os.Getenv(key)
}

var RedisClient *redis.Client

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})

	// Test connection
	_, err := RedisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
}
