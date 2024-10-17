package db

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func getEnvAsInt(name string, defaultVal int) int {
	valueStr := os.Getenv(name)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
}

func ConnectDB() (*gorm.DB, error) {
	// Load .env file only in development
	if os.Getenv("ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			fmt.Print("Error loading .env file")
		}
	}

	dbhost := os.Getenv("DB_HOST")
	dbport := getEnvAsInt("DB_PORT", 5432)
	dbuser := os.Getenv("DB_USERNAME")
	dbpassword := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	fmt.Printf("Connecting to DB at %s:%d with user %s\n", dbhost, dbport, dbuser)

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=require TimeZone=Asia/Jakarta", dbhost, dbuser, dbpassword, dbname, dbport)

	// Open the database connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: true, // Enable prepared statement caching
	})
	if err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
		return nil, err
	}

	// Set connection pool settings
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}
