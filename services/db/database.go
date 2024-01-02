package db

import (
	"fmt"
	"szimeth-jozef/clockmaster/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Connect(path string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("Error opening database: %w", err)
	}

	err = db.AutoMigrate(
		&models.WorkItem{},
		&models.WorkDay{},
	)
	if err != nil {
		return nil, fmt.Errorf("Error migrating database: %w", err)
	}

	return db, nil
}
