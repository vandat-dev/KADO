package initialize

import (
	"base_go_be/global"
	"base_go_be/pkg/setting"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

func LoadConfig() {
	// Load .env file
	err := godotenv.Load(".env")
	if err != nil {
		// If .env doesn't exist, try example.env
		err = godotenv.Load("example.env")
		if err != nil {
			fmt.Println("Warning: No .env or example.env file found, using system environment variables")
		}
	}

	// Load configuration from environment variables
	if err := loadConfigFromEnv(&global.Config); err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		panic(err)
	}
}

func loadConfigFromEnv(config *setting.Config) error {
	// Load Server settings
	config.Server = setting.ServerSetting{
		Port: getEnvAsInt("SERVER_PORT", 8082),
		Mode: getEnv("SERVER_MODE", "dev"),
	}

	// Load MySQL settings
	config.Mysql = setting.MySQLSetting{
		Host:            getEnv("MYSQL_HOST", "127.0.0.1"),
		Port:            getEnvAsInt("MYSQL_PORT", 33306),
		UserName:        getEnv("MYSQL_USERNAME", "admin"),
		Password:        getEnv("MYSQL_PASSWORD", "123123"),
		DBName:          getEnv("MYSQL_DATABASE", "go_database"),
		MaxIdleConn:     getEnvAsInt("MYSQL_MAX_IDLE_CONN", 10),
		MaxOpenConn:     getEnvAsInt("MYSQL_MAX_OPEN_CONN", 100),
		ConnMaxLifeTime: getEnvAsInt("MYSQL_CONN_MAX_LIFETIME", 3600),
	}

	// Load PostgreSQL settings
	config.Postgres = setting.PostgresSetting{
		Host:            getEnv("POSTGRES_HOST", "127.0.0.1"),
		Port:            getEnvAsInt("POSTGRES_PORT", 5432),
		UserName:        getEnv("POSTGRES_USERNAME", "postgres"),
		Password:        getEnv("POSTGRES_PASSWORD", "123123"),
		DBName:          getEnv("POSTGRES_DATABASE", "go_database"),
		SSLMode:         getEnv("POSTGRES_SSLMODE", "disable"),
		MaxIdleConn:     getEnvAsInt("POSTGRES_MAX_IDLE_CONN", 10),
		MaxOpenConn:     getEnvAsInt("POSTGRES_MAX_OPEN_CONN", 100),
		ConnMaxLifeTime: getEnvAsInt("POSTGRES_CONN_MAX_LIFETIME", 3600),
	}

	// Load Logger settings
	config.Logger = setting.LoggerSetting{
		LogLevel:    getEnv("LOG_LEVEL", "debug"),
		FileLogName: getEnv("LOG_FILE_NAME", "./storages/logs/development.xxx.log"),
		MaxBackups:  getEnvAsInt("LOG_MAX_BACKUPS", 3),
		MaxSize:     getEnvAsInt("LOG_MAX_SIZE", 500),
		MaxAge:      getEnvAsInt("LOG_MAX_AGE", 28),
		Compress:    getEnvAsBool("LOG_COMPRESS", true),
	}

	// Load Redis settings
	config.Redis = setting.RedisSetting{
		Host:     getEnv("REDIS_HOST", "redis"),
		Port:     getEnvAsInt("REDIS_PORT", 6379),
		Password: getEnv("REDIS_PASSWORD", ""),
		Database: getEnvAsInt("REDIS_DATABASE", 0),
	}

	// Load JWT settings
	config.JWT = setting.JWTSetting{
		SecretKey:     getEnv("JWT_SECRET_KEY", "your-secret-key-change-in-production"),
		TokenExpiry:   getEnvAsDuration("JWT_TOKEN_EXPIRY", 24*time.Hour),
		RefreshExpiry: getEnvAsDuration("JWT_REFRESH_EXPIRY", 7*24*time.Hour),
	}

	return nil
}

// Helper functions
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
}

func getEnvAsBool(name string, defaultVal bool) bool {
	valStr := getEnv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}
	return defaultVal
}

func getEnvAsDuration(name string, defaultVal time.Duration) time.Duration {
	valStr := getEnv(name, "")
	if val, err := time.ParseDuration(valStr); err == nil {
		return val
	}
	return defaultVal
}
