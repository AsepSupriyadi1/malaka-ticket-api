package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type AppConfig struct {
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
	JWTSecret  string
}

var App AppConfig

func LoadConfig() AppConfig {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found, reading config from environment")
	}

	App = AppConfig{
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBName:     os.Getenv("DB_NAME"),
		JWTSecret:  os.Getenv("JWT_SECRET"),
	}

	return App
}

func ConnectDatabase(cfg AppConfig) *gorm.DB {
	var db *gorm.DB
	var err error

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	maxAttempts := 10
	for i := 1; i <= maxAttempts; i++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			log.Println("✅ Connected to database!")
			break
		}

		log.Printf("⏳ Attempt %d: failed to connect to DB. Retrying in 3 seconds... Error: %v", i, err)
		time.Sleep(3 * time.Second)
	}

	if err != nil {
		log.Fatalf("❌ Could not connect to DB after %d attempts. Error: %v", maxAttempts, err)
	}

	return db
}
