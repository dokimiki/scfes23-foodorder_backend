package models

import (
	"time"
)

type Store struct {
	ID        uint32    `gorm:"primaryKey;auto_increment"`
	StrID     string    `gorm:"unique;not null"`
	Name      string    `gorm:"not null"`
	Location  string    `gorm:"not null"`
	Features  string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

type Menu struct {
	ID         uint32    `gorm:"primaryKey;auto_increment"`
	StrID      string    `gorm:"not null"`
	StoreStrID string    `gorm:"not null"`
	Name       string    `gorm:"not null"`
	Features   string    `gorm:"not null"`
	ImgURL     string    `gorm:"not null"`
	TimeToMake string    `gorm:"not null"`
	CreatedAt  time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

type MenuAllergen struct {
	MenuID    uint32    `gorm:"primaryKey"`
	Ebi       string    `gorm:"not null"`
	Kani      string    `gorm:"not null"`
	Komugi    string    `gorm:"not null"`
	Kurumi    string    `gorm:"not null"`
	Milk      string    `gorm:"not null"`
	Peanut    string    `gorm:"not null"`
	Soba      string    `gorm:"not null"`
	Tamago    string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

type StoreKeeper struct {
	ID         uint32    `gorm:"primaryKey;auto_increment"`
	Token      string    `gorm:"unique;not null"`
	IsApproved bool      `gorm:"not null;default:false"`
	StoreStrID string    `gorm:"not null"`
	Permission string    `gorm:"not null"`
	CreatedAt  time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

type User struct {
	ID          uint32    `gorm:"primaryKey;auto_increment"`
	Token       string    `gorm:"unique;not null"`
	CancelCount int       `gorm:"not null;default:0"`
	CreatedAt   time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

type Order struct {
	ID               uint32     `gorm:"primaryKey;auto_increment"`
	UserToken        string     `gorm:"not null"`
	StoreStrID       string     `gorm:"not null"`
	OrderStatus      string     `gorm:"not null;default:'choosing'"`
	IsCanceled       bool       `gorm:"not null;default:false"`
	IsPaid           bool       `gorm:"not null;default:false"`
	TimeOfCompletion *time.Time `gorm:"default:null"`
	CreatedAt        time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP"`
	DeletedAt        *time.Time `gorm:"default:null"`
}

type OrderItem struct {
	ID        uint32    `gorm:"primaryKey;auto_increment"`
	OrderID   uint32    `gorm:"not null"`
	MenuID    uint32    `gorm:"not null"`
	Quantity  int       `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

type StoreKeeperRequest struct {
	ID            uint32     `gorm:"primaryKey;auto_increment"`
	StoreKeeperID uint32     `gorm:"not null"`
	RequestToken  string     `gorm:"unique;not null"`
	CreatedAt     time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP"`
	DeletedAt     *time.Time `gorm:"default:null"`
}

type MenuDetail struct {
	MenuID      uint32    `gorm:"primaryKey"`
	Remaining   int       `gorm:"not null"`
	TicketPrice int       `gorm:"not null"`
	Discount    int       `gorm:"not null;default:0"`
	CreatedAt   time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time `gorm:"not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}

type Ticket struct {
	ID                uint32 `gorm:"primaryKey;auto_increment"`
	YenPricePerTicket int    `gorm:"not null"`
}

type Barcode struct {
	OrderID     uint32    `gorm:"primaryKey"`
	BarcodeData string    `gorm:"unique;not null"`
	CreatedAt   time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}
