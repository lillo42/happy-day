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
		ID         uint               `gorm:"primaryKey"`
		ExternalID uuid.UUID          `gorm:"uniqueIndex;not null"`
		Name       string             `gorm:"size:255;not null"`
		Price      float64            `gorm:"precision:10;scale:2;not null"`
		Discounts  []DiscountProducts `gorm:"foreignKey:ProductID;references:ID"`
		CreateAt   time.Time          `gorm:"not null"`
		UpdateAt   time.Time          `gorm:"not null"`
		Version    uint               `gorm:"not null"`
	}

	Discount struct {
		ID         uint               `gorm:"primaryKey"`
		ExternalID uuid.UUID          `gorm:"uniqueIndex;not null"`
		Name       string             `gorm:"size:255;not null"`
		Price      float64            `gorm:"precision:10;scale:2;not null"`
		Products   []DiscountProducts `gorm:"foreignKey:DiscountID;references:ID"`
		CreateAt   time.Time          `gorm:"not null"`
		UpdateAt   time.Time          `gorm:"not null"`
		Version    uint               `gorm:"not null"`
	}

	DiscountProducts struct {
		DiscountID uint     `gorm:"primaryKey;autoIncrement:false;not null"`
		Discount   Discount `gorm:"foreignKey:DiscountID;references:ID"`
		ProductID  uint     `gorm:"primaryKey;autoIncrement:false;not null"`
		Product    Product  `gorm:"foreignKey:ProductID;references:ID"`
		Quantity   uint     `gorm:"not null"`
	}
)
