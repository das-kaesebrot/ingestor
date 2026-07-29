package repository

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"dev.kaesebrot.eu/go/ingestor/internal/utility"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func New(dbFile string) (*Repository, error) {
	slog.Info("Reading SQLite database file", "dbFile", dbFile)
	err := utility.CheckFileAccess(dbFile)

	if errors.Is(err, os.ErrNotExist) {
		slog.Info("database file doesn't exist yet, creating it", "dbFile", dbFile)
	} else if err != nil {
		return nil, fmt.Errorf("failed reading database file: %w", err)
	}

	db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	db.AutoMigrate(&Project{}, &UploadUser{})

	return &Repository{db: db}, nil
}
