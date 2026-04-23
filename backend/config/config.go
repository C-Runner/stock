package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"

	"backend/models"
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

// MigrateDB runs database migrations for user_id fields
func MigrateDB() error {
	// === Stocks table (user_id as part of composite PK) ===
	if err := DB.Exec("ALTER TABLE stocks ADD COLUMN IF NOT EXISTS user_id VARCHAR(255)").Error; err != nil {
		return fmt.Errorf("failed to add user_id to stocks: %w", err)
	}

	if err := DB.Exec("ALTER TABLE stocks DROP CONSTRAINT IF EXISTS stocks_pkey").Error; err != nil {
		log.Printf("Note: stocks_pkey may not exist or already dropped: %v", err)
	}

	if err := DB.Exec("ALTER TABLE stocks ADD PRIMARY KEY (code, user_id)").Error; err != nil {
		log.Printf("Note: stocks composite pk may already exist: %v", err)
	}

	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_stocks_user_id ON stocks(user_id)").Error; err != nil {
		return fmt.Errorf("failed to create idx_stocks_user_id: %w", err)
	}

	// === Watchlist table (no user_id, keep as backup data table) ===
	// Ensure watchlist table has proper primary key (code only)
	if err := DB.Exec("ALTER TABLE watchlist DROP CONSTRAINT IF EXISTS watchlist_pkey").Error; err != nil {
		log.Printf("Note: watchlist_pkey may not exist: %v", err)
	}

	if err := DB.Exec("ALTER TABLE watchlist ADD PRIMARY KEY (code)").Error; err != nil {
		log.Printf("Note: watchlist pk may already exist: %v", err)
	}

	// === UserWatchlist table (new user-watchlist association table) ===
	// Create user_watchlist table
	if err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS user_watchlist (
			user_id VARCHAR(255) NOT NULL,
			code VARCHAR(20) NOT NULL,
			added_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (user_id, code)
		)
	`).Error; err != nil {
		return fmt.Errorf("failed to create user_watchlist table: %w", err)
	}

	// Create index on user_id
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_user_watchlist_user_id ON user_watchlist(user_id)").Error; err != nil {
		return fmt.Errorf("failed to create idx_user_watchlist_user_id: %w", err)
	}

	log.Println("Database migration completed successfully")
	return nil
}

// SeedDefaultUser creates default admin user if not exists
func SeedDefaultUser() error {
	var count int64
	if err := DB.Model(&models.User{}).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to count users: %w", err)
	}

	if count > 0 {
		log.Println("Users already exist, skipping seed")
		return nil
	}

	cfg := LoadConfig()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(cfg.LoginPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	adminUser := models.User{
		Username: cfg.LoginUsername,
		Password: string(hashedPassword),
	}

	if err := DB.Create(&adminUser).Error; err != nil {
		return fmt.Errorf("failed to create default user: %w", err)
	}

	log.Printf("Default user created: %s", cfg.LoginUsername)
	return nil
}
