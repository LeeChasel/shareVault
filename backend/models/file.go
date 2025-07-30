package models

import (
	"time"

	"github.com/google/uuid"
)

type File struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index:idx_user_filehash,unique"`
	FileName  string    `gorm:"not null"`
	FilePath  string    `gorm:"not null"`
	FileHash  string    `gorm:"not null;index:idx_user_filehash,unique"`
	FileSize  int64     `gorm:"not null"`
	MimeType  string    `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}