package models

import (
	"time"
)

type Menu struct {
	ID     uint32 `gorm:"primaryKey"`
	Name   string `gorm:"not null"`
	Price  int    `gorm:"not null"`
	ImgURL string `gorm:"not null"`
}

type MenuAllergen struct {
	MenuID uint32 `gorm:"primaryKey"`
	Ebi    string `gorm:"not null"`
	Kani   string `gorm:"not null"`
	Komugi string `gorm:"not null"`
	Kurumi string `gorm:"not null"`
	Milk   string `gorm:"not null"`
	Peanut string `gorm:"not null"`
	Soba   string `gorm:"not null"`
	Tamago string `gorm:"not null"`
}

type StoreKeeper struct {
	ID         uint32 `gorm:"primaryKey"`
	Token      string `gorm:"unique;not null"`
	IsApproved bool   `gorm:"not null;default:false"`
}

type User struct {
	ID        uint32    `gorm:"primaryKey"`
	Token     string    `gorm:"unique;not null"`
	CreatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

type MenuDetail struct {
	MenuID    uint32 `gorm:"primaryKey"`
	Remaining int    `gorm:"not null"`
	Discount  int    `gorm:"not null;default:0"`
}

type Order struct {
	ID               uint32 `gorm:"primaryKey"`
	UserID           uint32 `gorm:"not null"`
	OrderStatus      string `gorm:"not null;default:choosing"`
	NumberTag        int    `gorm:"not null;default:0"`
	IsCanceled       bool   `gorm:"not null;default:false"`
	IsPaid           bool   `gorm:"not null;default:false"`
	IsMobileOrder    bool   `gorm:"not null;default:false"`
	IsInvitationSent bool   `gorm:"not null;default:false"`
	TimeOfCompletion time.Time
	CreatedAt        time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

type OrderItem struct {
	ID       uint32 `gorm:"primaryKey"`
	OrderID  uint32 `gorm:"not null"`
	MenuID   uint32 `gorm:"not null"`
	Quantity int    `gorm:"not null"`
}

type Barcode struct {
	OrderID     uint32    `gorm:"primaryKey"`
	BarcodeData string    `gorm:"unique;not null"`
	CreatedAt   time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}
