package models

import (
	"time"
)

type File struct {
	FileID           string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"fileId"`
	FileUri          string    `gorm:"not null" json:"fileUri"`
	FileThumbnailUri string    `gorm:"not null" json:"fileThumbnailUri"`
	CreatedAt        time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt        time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updatedAt"`
}
