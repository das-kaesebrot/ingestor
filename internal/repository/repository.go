package repository

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"uuid"

	"dev.kaesebrot.eu/go/ingestor/internal/utility"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
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

	db.AutoMigrate(&Project{}, &UploadUser{}, &ImageFile{}, &VideoFile{})

	return &Repository{DB: db}, nil
}

func (r *Repository) Create(ctx context.Context, object any) error {
	return r.DB.WithContext(ctx).Create(object).Error
}

func (r *Repository) Save(ctx context.Context, object any) error {
	return r.DB.WithContext(ctx).Save(object).Error
}

func (r *Repository) ProjectByShareToken(ctx context.Context, shareToken string) (*Project, error) {
	var m Project
	err := r.DB.WithContext(ctx).Where(&Project{ShareToken: shareToken}).Take(&m).Error
	return &m, err
}

func (r *Repository) ProjectByAdminToken(ctx context.Context, adminToken string) (*Project, error) {
	var m Project
	err := r.DB.WithContext(ctx).Where(&Project{AdminToken: adminToken}).Take(&m).Error
	return &m, err
}

func (r *Repository) ProjectIdByShareToken(ctx context.Context, shareToken string) (uuid.UUID, error) {
	var m Project
	err := r.DB.WithContext(ctx).Select("id").Where(&Project{ShareToken: shareToken}).Take(&m).Error
	if err != nil {
		return uuid.UUID{}, err
	}
	return m.ID, nil
}

func (r *Repository) ProjectIdByAdminToken(ctx context.Context, adminToken string) (uuid.UUID, error) {
	var m Project
	err := r.DB.WithContext(ctx).Select("id").Where(&Project{AdminToken: adminToken}).Take(&m).Error
	if err != nil {
		return uuid.UUID{}, err
	}
	return m.ID, nil
}

func (r *Repository) UploadUserById(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) (*UploadUser, error) {
	var m UploadUser
	err := r.DB.WithContext(ctx).Where(&UploadUser{ID: userID, ProjectID: projectID}).Take(&m).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *Repository) UploadUsersByProjectWithLimitAndOffset(ctx context.Context, projectID uuid.UUID, limit int, offset int) ([]UploadUser, error) {
	var users []UploadUser
	err := r.DB.WithContext(ctx).Where(&UploadUser{ProjectID: projectID}).Order("created_at ASC").Limit(limit).Offset(offset).Find(&users).Error
	if err != nil {
		return []UploadUser{}, err
	}
	return users, nil
}

func (r *Repository) CountUploadUsersByProject(ctx context.Context, projectID uuid.UUID) (int64, error) {
	var count int64
	err := r.DB.WithContext(ctx).Model(&UploadUser{}).Select("id").Where(&UploadUser{ProjectID: projectID}).Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *Repository) DeleteProject(ctx context.Context, projectID uuid.UUID) error {
	result := r.DB.WithContext(ctx).Delete(&Project{ID: projectID})
	return result.Error
}

func (r *Repository) DeleteProjectByAdminToken(ctx context.Context, adminToken string) error {
	result := r.DB.Where(&Project{AdminToken: adminToken}).Delete(ctx)
	return result.Error
}

func (r *Repository) DeleteUploadUser(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) error {
	u, err := r.UploadUserById(ctx, projectID, userID)
	if err != nil {
		return err
	}

	result := r.DB.Delete(u)
	return result.Error
}
