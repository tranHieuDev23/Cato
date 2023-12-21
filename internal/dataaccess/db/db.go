package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitializeDB(connectionString string) (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(connectionString), &gorm.Config{})
}
