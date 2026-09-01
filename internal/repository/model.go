package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Project struct {
	gorm.Model
	ID          uuid.UUID `gorm:"primaryKey,type:uuid"`
	ShareKey    string    `gorm:"uniqueIndex"`
	AdminKey    string    `gorm:"uniqueIndex"`
	Description string
	UploadUsers []UploadUser
}

func (p *Project) BeforeCreate(tx *gorm.DB) (err error) {
	newUuid, err := uuid.NewV7()
	if err != nil {
		return err
	}
	p.ID = newUuid
	return
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
