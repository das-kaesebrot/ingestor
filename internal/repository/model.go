package repository

import (
	"errors"
	"time"
	"uuid"

	"dev.kaesebrot.eu/go/ingestor/internal/mediaprocessor"
	"dev.kaesebrot.eu/go/ingestor/internal/utility"
	"gorm.io/gorm"
)

const (
	defaultTokenLength   = 16
	defaultTokenAlphabet = utility.TokenAlphabetHumanReadable
)

func GenerateRandomRepoToken() (string, error) {
	return utility.GenerateRandomToken(defaultTokenLength, defaultTokenAlphabet)
}

type ShareTokenPermission int

const (
	SHARE_TOKEN_UPLOAD_ALLOWED ShareTokenPermission = iota
	SHARE_TOKEN_DOWNLOAD_ALLOWED
)

type Project struct {
	gorm.Model
	ID                    uuid.UUID `gorm:"primaryKey,type:uuid"`
	ShareToken            string    `gorm:"not null;uniqueIndex"`
	AdminToken            string    `gorm:"not null;uniqueIndex"`
	Description           string
	Start                 *time.Time
	End                   *time.Time
	ExpiresAt             time.Time
	ShareTokenPermissions []ShareTokenPermission `gorm:"serializer:json"`
	UploadUsers           []UploadUser
	MediaFiles            []MediaFileBase
}

func (p *Project) Validate() error {
	if p.Start != nil && p.End != nil && (p.Start.After(*p.End) || p.Start.Equal(*p.End)) {
		return errors.New("Project start can't be at or after project end!")
	}

	if p.ExpiresAt.Before(time.Now()) {
		return errors.New("Project expiry can't be in the past!")
	}

	return nil
}

func (p *Project) BeforeCreate(tx *gorm.DB) error {
	newUuid := uuid.NewV7()
	p.ID = newUuid

	newShareToken, err := GenerateRandomRepoToken()
	if err != nil {
		return err
	}
	p.ShareToken = newShareToken

	newAdminToken, err := GenerateRandomRepoToken()
	if err != nil {
		return err
	}
	p.AdminToken = newAdminToken

	return nil
}

type UploadUser struct {
	gorm.Model
	ID        uuid.UUID `gorm:"primaryKey,type:uuid"`
	Name      string
	ProjectID uuid.UUID
	MediaFile []MediaFileBase
}

func (u *UploadUser) BeforeCreate(tx *gorm.DB) error {
	newUuid := uuid.NewV7()
	u.ID = newUuid
	return nil
}

type MediaFileBase struct {
	gorm.Model
	ID           uuid.UUID `gorm:"primaryKey,type:uuid"`
	UploadUserID *uuid.UUID
	ProjectID    uuid.UUID
}

type ImageFile struct {
	MediaFileBase
	mediaprocessor.ImageMetadata
}

type VideoFile struct {
	MediaFileBase
	mediaprocessor.VideoMetadata
}

func (m *MediaFileBase) BeforeCreate(tx *gorm.DB) error {
	// don't generate a new uuid if we already have a value (this would be generated during file upload)
	if m.ID != uuid.Nil() {
		return nil
	}

	newUuid := uuid.NewV7()
	m.ID = newUuid
	return nil
}
