package db

import (
	"os"

	"github.com/Giankrp/AlcatrazBack/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DbConnection() error {

	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	DB = db
	return nil
}

func AutoMigrate() error {
	return DB.AutoMigrate(&models.User{}, &models.VaultItem{}, &models.VaultFolder{}, &models.Session{})
}
