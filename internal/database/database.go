package database

import (
	"os"

	"github.com/0xirvan/go-jwt-auth/internal/models"
	"github.com/0xirvan/go-jwt-auth/pkg/logger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Connect() *gorm.DB {
	databasePath := os.Getenv("DATABASE_PATH")
	if databasePath == "" {
		databasePath = "auth.db"
	}

	db, err := gorm.Open(sqlite.Open(databasePath), &gorm.Config{
		PrepareStmt:            true,
		TranslateError:         true,
		SkipDefaultTransaction: true,
	})

	if err != nil {
		logger.Log.Fatalf("failed to connect to database: %v", err)
	}

	// Migrate the models
	err = db.AutoMigrate(&models.User{}, &models.RefreshToken{})
	if err != nil {
		logger.Log.Fatalf("failed to migrate database: %v", err)
	}

	logger.Log.Info("Connected to database successfully")
	return db
}
