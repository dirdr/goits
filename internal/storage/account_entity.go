package storage

import "time"

type GormAccount struct {
	ID        uint      `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

func (GormAccount) TableName() string {
	return "accounts"
}
