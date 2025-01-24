package models

type ProductCategory struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Type string `gorm:"unique;not null" json:"type"`
}
