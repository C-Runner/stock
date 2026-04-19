package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	PostgresHost         string
	PostgresUser         string
	PostgresPassword     string
	PostgresDBName       string
	PostgresPort         string
	ServerPort           string
	MaxIdleConns         int
	MaxOpenConns         int
	ConnMaxLifetimeHours int
	ConnectTimeout       int
	MaxRetries           int
	LoginUsername        string
	LoginPassword        string
	JWTSecret            string
}

var DB *gorm.DB

func LoadConfig() *Config {
	return &Config{
		PostgresHost:         getEnv("POSTGRES_HOST", "43.139.206.109"),
		PostgresUser:         getEnv("POSTGRES_USER", "username"),
		PostgresPassword:     getEnv("POSTGRES_PASSWORD", "your_password"),
		PostgresDBName:       getEnv("POSTGRES_DB", "stock"),
		PostgresPort:         getEnv("POSTGRES_PORT", "5432"),
		ServerPort:           getEnv("SERVER_PORT", "8080"),
		MaxIdleConns:         getEnvInt("POSTGRES_MAX_IDLE_CONNS", 10),
		MaxOpenConns:         getEnvInt("POSTGRES_MAX_OPEN_CONNS", 100),
		ConnMaxLifetimeHours: getEnvInt("POSTGRES_CONN_MAX_LIFETIME_HOURS", 1),
		ConnectTimeout:       getEnvInt("POSTGRES_CONNECT_TIMEOUT", 10),
		MaxRetries:           getEnvInt("POSTGRES_MAX_RETRIES", 3),
		LoginUsername:        getEnv("LOGIN_USERNAME", "admin"),
		LoginPassword:        getEnv("LOGIN_PASSWORD", "admin123"),
		JWTSecret:            getEnv("JWT_SECRET", "stock-app-secret-key"),
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		if intVal, err := strconv.Atoi(val); err == nil {
			return intVal
		}
	}
	return defaultVal
}

func ConnectPostgres(cfg *Config) error {
	var lastErr error
	for i := 0; i < cfg.MaxRetries; i++ {
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable connect_timeout=%d",
			cfg.PostgresHost,
			cfg.PostgresUser,
			cfg.PostgresPassword,
			cfg.PostgresDBName,
			cfg.PostgresPort,
			cfg.ConnectTimeout,
		)

		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			sqlDB, err := db.DB()
			if err != nil {
				return fmt.Errorf("failed to get underlying sql.DB: %w", err)
			}
			sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
			sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
			sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetimeHours) * time.Hour)

			DB = db
			log.Println("Connected to PostgreSQL successfully!")
			return nil
		}
		lastErr = err
		log.Printf("PostgreSQL connection attempt %d/%d failed: %v", i+1, cfg.MaxRetries, err)
		if i < cfg.MaxRetries-1 {
			time.Sleep(time.Duration(2<<i) * time.Second)
		}
	}
	return fmt.Errorf("failed to connect to PostgreSQL after %d attempts: %w", cfg.MaxRetries, lastErr)
}

func PingPostgres() error {
	if DB == nil {
		return fmt.Errorf("database connection not initialized")
	}
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

func DisconnectPostgres() {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err == nil {
			sqlDB.Close()
			log.Println("Disconnected from PostgreSQL")
		}
	}
}
