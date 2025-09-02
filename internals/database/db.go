package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/amarjeet-choudhary666/ai_resume_screener/internals/config"
	"github.com/amarjeet-choudhary666/ai_resume_screener/internals/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB(cfg *config.Config) (*gorm.DB, error) {
	dbUrl := cfg.Database

	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := sqlDB.PingContext(context.Background()); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS pgcrypto").Error; err != nil {
		log.Fatalf("Failed to initialize pgcrypto: %v", err)
	}

	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("Failed to migrate user: %v", err)
	}

	DB = db

	fmt.Println("✅Successfully connected to database")
	return db, nil
}
