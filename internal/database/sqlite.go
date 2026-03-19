package database

import (
	"gogo/internal/pet"
	"os/user"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitSQLite() (*gorm.DB, error) {
	db, err := gorm.Open(
		sqlite.Open("internal/database/gogo.db?_foreign_keys=on"),
		&gorm.Config{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&user.User{},
		&pet.Pet{},
	)
	if err != nil {
		return err
	}
	return nil
}
