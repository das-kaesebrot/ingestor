package repository

import (
	"errors"
	"time"

	"github.com/google/uuid"

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
	newUuid, err := uuid.NewV7()
	if err != nil {
		return err
	}
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
	ID         uuid.UUID `gorm:"primaryKey,type:uuid"`
	Name       string
	ProjectID  uuid.UUID
	ImageFiles []ImageFile
	VideoFiles []VideoFile
}

func (u *UploadUser) BeforeCreate(tx *gorm.DB) error {
	newUuid, err := uuid.NewV7()
	if err != nil {
		return err
	}
	u.ID = newUuid
	return nil
}

type MediaFile struct {
	gorm.Model
	ID                      uuid.UUID `gorm:"primaryKey,type:uuid"`
	UploadUserID            uuid.UUID
	ProjectID               uuid.UUID
	OriginalFilename        string
	CaptureTime             time.Time
	CaptureCorrectionOffset time.Duration
}

func (m *MediaFile) BeforeCreate(tx *gorm.DB) error {
	newUuid, err := uuid.NewV7()
	if err != nil {
		return err
	}
	m.ID = newUuid
	return nil
}

type VideoFile struct {
	MediaFile
	VideoMetadata
}

type ImageFile struct {
	MediaFile
	ImageMetadata
}

type ImageMetadata struct {
	Width     int64
	Height    int64
	Codec     string
	Container string
}

type VideoMetadata struct {
	Width     int64
	Height    int64
	Duration  time.Duration
	Codec     string
	Container string
}
