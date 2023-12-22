package db

import (
	"errors"
	"fmt"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/tranHieuDev23/cato/internal/configs"
)

func InitializeDB(
	databaseConfig configs.Database,
) (*gorm.DB, error) {
	databaseFilePath, err := filepath.Abs(databaseConfig.FilePath)
	if err != nil {
		return nil, errors.New("failed to get abs database file path")
	}

	return gorm.Open(
		sqlite.Open(fmt.Sprintf("file://%s", databaseFilePath)),
		&gorm.Config{},
	)
}
