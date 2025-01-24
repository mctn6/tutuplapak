package models

import "time"

type Product struct {
	ID        int       `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:32;not null" json:"name"`
	Category  string    `gorm:"not null" json:"category"`
	Qty       int       `gorm:"not null;check:qty >= 1" json:"qty"`
	Price     float64   `gorm:"not null;check:price >= 100" json:"price"`
	SKU       string    `gorm:"size:32;not null" json:"sku"`
	FileID    string    `gorm:"type:uuid;not null" json:"fileId"`
	File      File      `gorm:"foreignKey:FileID" json:"file"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updatedAt"`
}
