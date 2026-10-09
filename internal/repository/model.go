package repository

import (
	"time"

	"github.com/google/uuid"

	"dev.kaesebrot.eu/go/ingestor/internal/utility"
	"gorm.io/gorm"
)

type Project struct {
	gorm.Model
	ID          uuid.UUID `gorm:"primaryKey,type:uuid"`
	ShareToken  string    `gorm:"not null;uniqueIndex"`
	AdminToken  string    `gorm:"not null;uniqueIndex"`
	Description string
	UploadUsers []UploadUser
}

func (p *Project) BeforeCreate(tx *gorm.DB) error {
	newUuid, err := uuid.NewV7()
	if err != nil {
		return err
	}
	p.ID = newUuid

	defaultTokenLength := 16
	defaultTokenAlphabet := utility.TokenAlphabetHumanReadable

	newShareToken, err := utility.GenerateRandomToken(defaultTokenLength, defaultTokenAlphabet)
	if err != nil {
		return err
	}
	p.ShareToken = newShareToken

	newAdminToken, err := utility.GenerateRandomToken(defaultTokenLength, defaultTokenAlphabet)
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
