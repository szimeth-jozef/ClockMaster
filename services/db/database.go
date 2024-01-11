package db

import (
	"fmt"
	"os"
	"path"
	"szimeth-jozef/clockmaster/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Connect(dbFolderPath *string) (*gorm.DB, error) {
	const DB_FILE_NAME = "db.sqlite3"
	var dbDirectory string
	if dbFolderPath == nil || *dbFolderPath == "" {

		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("Error getting user home directory: %w", err)
		}

		dbDirectory = path.Join(
			homeDir,
			".clockmaster",
		)
	} else {
		dbDirectory = *dbFolderPath
	}

	err := os.MkdirAll(dbDirectory, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("Error creating database directory: %w", err)
	}

	dbFilePath := path.Join(dbDirectory, DB_FILE_NAME)

	db, err := gorm.Open(sqlite.Open(dbFilePath), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("Error opening database: %w", err)
	}

	fmt.Println("Database location at:", dbFilePath)

	err = db.AutoMigrate(
		&models.WorkItem{},
		&models.WorkDay{},
	)
	if err != nil {
		return nil, fmt.Errorf("Error migrating database: %w", err)
	}

	return db, nil
}
