package db

import (
	"fmt"
	"os"
	"strconv"

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

	dbhost :="aws-0-eu-central-1.pooler.supabase.com"
	dbport := 6543
	dbuser := "postgres.ndlncbadozimymlhaeyw"
	dbpassword := "AsllanPireva69!Nice"
	dbname := "postgres"


	// Adjust sslmode as needed
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=require TimeZone=Asia/Jakarta", dbhost, dbuser, dbpassword, dbname, dbport)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
