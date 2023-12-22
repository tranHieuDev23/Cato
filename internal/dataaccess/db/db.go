package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitializeDB() (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(""), &gorm.Config{})
}
