package infra

import (
	"github.com/google/uuid"
	"time"
)

type (
	Customer struct {
		ID         uint      `gorm:"primaryKey"`
		ExternalID uuid.UUID `gorm:"uniqueIndex;not null"`
		Name       string    `gorm:"size:255;not null"`
		Comment    string
		Phones     string
		Pix        string    `gorm:"size:50"`
		CreateAt   time.Time `gorm:"not null"`
		UpdateAt   time.Time `gorm:"not null"`
		Version    uint      `gorm:"not null"`
	}

	Product struct {
		ID         uint         `gorm:"primaryKey"`
		ExternalID uuid.UUID    `gorm:"uniqueIndex;not null"`
		Name       string       `gorm:"size:255;not null"`
		Price      float64      `gorm:"precision:10;scale:2;not null"`
		Products   []BoxProduct `gorm:"foreignKey:ParentID"`
		CreateAt   time.Time    `gorm:"not null"`
		UpdateAt   time.Time    `gorm:"not null"`
		Version    uint         `gorm:"not null"`
	}

	BoxProduct struct {
		ParentID  uint    `gorm:"primaryKey;autoIncrement:false;not null"`
		Parent    Product `gorm:"foreignKey:ParentID"`
		ProductID uint    `gorm:"primaryKey;autoIncrement:false;not null"`
		Product   Product `gorm:"foreignKey:ProductID"`
		Quantity  uint    `gorm:"not null"`
	}
)
