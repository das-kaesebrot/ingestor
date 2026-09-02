package repository

import (
	"dev.kaesebrot.eu/go/ingestor/internal/utility"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Project struct {
	gorm.Model
	ID          uuid.UUID `gorm:"primaryKey,type:uuid"`
	ShareToken  string    `gorm:"unique"`
	AdminToken  string    `gorm:"unique"`
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
	ID        uuid.UUID `gorm:"primaryKey,type:uuid"`
	Name      string
	ProjectID uuid.UUID
}

func (u *UploadUser) BeforeCreate(tx *gorm.DB) (err error) {
	newUuid, err := uuid.NewV7()
	if err != nil {
		return err
	}
	u.ID = newUuid
	return
}

type ImageFile struct {
	gorm.Model
	UserID         uuid.UUID
	FilepathOnDisk string
	ImageMetadata
}

type ImageMetadata struct {
	Width     int64
	Height    int64
	Codec     string
	Container string
}
