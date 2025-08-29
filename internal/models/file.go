package models

import (
	"time"

	"github.com/google/uuid"
)

type File struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"userId"`
	FileName  string    `gorm:"not null" json:"fileName"`
	FilePath  string    `gorm:"not null" json:"filePath"`
	FileHash  string    `gorm:"not null" json:"fileHash"`
	FileSize  int64     `gorm:"not null" json:"fileSize"`
	MimeType  string    `gorm:"not null" json:"mimeType"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
