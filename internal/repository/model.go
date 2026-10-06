package repository

import (
	"time"
	stdlibuuid "uuid"

	"dev.kaesebrot.eu/go/ingestor/internal/utility"
	"gorm.io/gorm"
)

type Project struct {
	gorm.Model
	ID          stdlibuuid.UUID `gorm:"primaryKey,type:uuid"`
	ShareToken  string          `gorm:"unique"`
	AdminToken  string          `gorm:"unique"`
	Description string
	UploadUsers []UploadUser
}

func (p *Project) BeforeCreate(tx *gorm.DB) error {
	newUuid := stdlibuuid.NewV7()
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
	ID         stdlibuuid.UUID `gorm:"primaryKey,type:uuid"`
	Name       string
	ProjectID  stdlibuuid.UUID
	ImageFiles []ImageFile
	VideoFiles []VideoFile
}

func (u *UploadUser) BeforeCreate(tx *gorm.DB) (err error) {
	newUuid := stdlibuuid.NewV7()
	u.ID = newUuid
	return
}

type MediaFile struct {
	gorm.Model
	ID                      stdlibuuid.UUID `gorm:"primaryKey,type:uuid"`
	UploadUserID            stdlibuuid.UUID
	ProjectID               stdlibuuid.UUID
	OriginalFilename        string
	CaptureTime             time.Time
	CaptureCorrectionOffset time.Duration
}

func (m *MediaFile) BeforeCreate(tx *gorm.DB) (err error) {
	newUuid := stdlibuuid.NewV7()
	m.ID = newUuid
	return
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
