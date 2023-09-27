package models

import (
	"time"
)

type Menu struct {
	ID     uint32 `gorm:"primary_key;AUTO_INCREMENT"`
	Name   string `gorm:"not null"`
	Price  int    `gorm:"not null"`
	ImgUrl string `gorm:"not null"`
	IsShow bool   `gorm:"not null;default:true"`
}

type MenuAllergen struct {
	MenuID uint32 `gorm:"primary_key"`
	Ebi    string `gorm:"not null"`
	Kani   string `gorm:"not null"`
	Komugi string `gorm:"not null"`
	Kurumi string `gorm:"not null"`
	Milk   string `gorm:"not null"`
	Peanut string `gorm:"not null"`
	Soba   string `gorm:"not null"`
	Tamago string `gorm:"not null"`
}

type User struct {
	ID           uint32    `gorm:"primary_key;AUTO_INCREMENT"`
	Token        string    `gorm:"unique;not null"`
	IsInvitation bool      `gorm:"not null;default:false"`
	IsOrdered    bool      `gorm:"not null;default:false"`
	InviteCoupon string    `gorm:"not null;default:'none'"`
	BulkCoupon   string    `gorm:"not null;default:'none'"`
	CreatedAt    time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

type ExceptionPotatoOrder struct {
	ID       uint32 `gorm:"primary_key"`
	Quantity int    `gorm:"not null;default:0"`
}

type Order struct {
	ID               uint32 `gorm:"primary_key;AUTO_INCREMENT"`
	UserID           uint32 `gorm:"not null"`
	OrderStatus      string `gorm:"not null;default:'choosing'"`
	NumberTag        int    `gorm:"not null;default:0"`
	IsPaid           bool   `gorm:"not null;default:false"`
	IsMobileOrder    bool   `gorm:"not null;default:false"`
	TimeOfCompletion time.Time
	CreatedAt        time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

type OrderItem struct {
	ID       uint32 `gorm:"primary_key;AUTO_INCREMENT"`
	OrderID  uint32 `gorm:"not null"`
	MenuID   uint32 `gorm:"not null"`
	Quantity int    `gorm:"not null"`
}

type Barcode struct {
	OrderID     uint32    `gorm:"primary_key"`
	BarcodeData string    `gorm:"unique;not null"`
	CreatedAt   time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

type PotatoDetail struct {
	ID        uint32 `gorm:"primary_key"`
	Remaining int    `gorm:"not null"`
}
