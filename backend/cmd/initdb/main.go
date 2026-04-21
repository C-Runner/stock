package main

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=43.139.206.109 user=username password=your_password dbname=stock port=5432 sslmode=disable"
	// 先连接到默认的 postgres 数据库来创建 stock 数据库
	dsnDefault := "host=43.139.206.109 user=username password=your_password dbname=postgres port=5432 sslmode=disable"
	dbDefault, err := gorm.Open(postgres.Open(dsnDefault), &gorm.Config{})
	if err != nil {
		log.Fatal("连接 postgres 数据库失败:", err)
	}

	// 创建 stock 数据库（如果不存在）
	dbDefault.Exec("CREATE DATABASE stock")

	sqlDB, _ := dbDefault.DB()
	sqlDB.Close()

	// 连接到 stock 数据库
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("连接 stock 数据库失败:", err)
	}

	sqls := []string{
		`CREATE TABLE IF NOT EXISTS stocks (
			code VARCHAR(20) PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			current_price DECIMAL(15, 2) NOT NULL DEFAULT 0,
			quantity INT NOT NULL DEFAULT 0,
			buy_price DECIMAL(15, 2) NOT NULL DEFAULT 0,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS watchlist (
			code VARCHAR(20) PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			added_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE OR REPLACE FUNCTION update_updated_at_column()
		RETURNS TRIGGER AS $$
		BEGIN
			NEW.updated_at = CURRENT_TIMESTAMP;
			RETURN NEW;
		END;
		$$ language 'plpgsql'`,
		`DROP TRIGGER IF EXISTS update_stocks_updated_at ON stocks`,
		`CREATE TRIGGER update_stocks_updated_at BEFORE UPDATE ON stocks
		FOR EACH ROW EXECUTE FUNCTION update_updated_at_column()`,

		// stock_daily_snapshots table for backup
		`CREATE TABLE IF NOT EXISTS stock_daily_snapshots (
			code VARCHAR(20) NOT NULL,
			date VARCHAR(20) NOT NULL,
			name VARCHAR(100),
			open DECIMAL(15, 2) DEFAULT 0,
			high DECIMAL(15, 2) DEFAULT 0,
			low DECIMAL(15, 2) DEFAULT 0,
			close DECIMAL(15, 2) DEFAULT 0,
			volume BIGINT DEFAULT 0,
			amount DECIMAL(20, 2) DEFAULT 0,
			turnover_rate DECIMAL(10, 4) DEFAULT 0,
			ma5 DECIMAL(15, 2) DEFAULT 0,
			ma10 DECIMAL(15, 2) DEFAULT 0,
			ma20 DECIMAL(15, 2) DEFAULT 0,
			ma60 DECIMAL(15, 2) DEFAULT 0,
			ema12 DECIMAL(15, 2) DEFAULT 0,
			ema26 DECIMAL(15, 2) DEFAULT 0,
			rsi6 DECIMAL(10, 4) DEFAULT 0,
			rsi12 DECIMAL(10, 4) DEFAULT 0,
			rsi24 DECIMAL(10, 4) DEFAULT 0,
			dif DECIMAL(15, 6) DEFAULT 0,
			dea DECIMAL(15, 6) DEFAULT 0,
			macd DECIMAL(15, 6) DEFAULT 0,
			kdjk DECIMAL(10, 4) DEFAULT 0,
			kdjd DECIMAL(10, 4) DEFAULT 0,
			kdjj DECIMAL(10, 4) DEFAULT 0,
			boll_upper DECIMAL(15, 2) DEFAULT 0,
			boll_mid DECIMAL(15, 2) DEFAULT 0,
			boll_lower DECIMAL(15, 2) DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (code, date)
		)`,
	}

	for _, sql := range sqls {
		if err := db.Exec(sql).Error; err != nil {
			log.Fatal("执行 SQL 失败:", err)
		}
		fmt.Println("执行成功:", sql[:30], "...")
	}

	// 验证表是否创建成功
	var tables []string
	db.Raw("SELECT tablename FROM pg_tables WHERE schemaname = 'public'").Scan(&tables)
	fmt.Println("当前表:", tables)
}
