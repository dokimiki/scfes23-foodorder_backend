package models

import (
	"time"
)

type Store struct {
	ID        uint32 `gorm:"primary_key;auto_increment"`
	StrID     string `gorm:"unique;not null"`
	Name      string `gorm:"not null"`
	Location  string `gorm:"not null"`
	Features  string `gorm:"not null"`
	CreatedAt time.Time
}

type StoreKeeper struct {
	ID         uint32 `gorm:"primary_key;auto_increment"`
	Token      string `gorm:"unique;not null"`
	IsApproved bool   `gorm:"not null"`
	StoreID    uint32 `gorm:"not null"`
	Permission string `gorm:"not null"`
	CreatedAt  time.Time
}

type User struct {
	ID          uint32 `gorm:"primary_key;auto_increment"`
	Token       string `gorm:"unique;not null"`
	CancelCount int    `gorm:"not null"`
	CreatedAt   time.Time
}

type Order struct {
	ID               uint32 `gorm:"primary_key;auto_increment"`
	UserID           uint32 `gorm:"not null"`
	StoreID          uint32 `gorm:"not null"`
	OrderStatus      string `gorm:"not null"`
	IsCanceled       bool   `gorm:"not null"`
	IsPaid           bool   `gorm:"not null"`
	TimeOfCompletion time.Time
	CreatedAt        time.Time
}

type Menu struct {
	ID         uint32 `gorm:"primary_key;auto_increment"`
	StrID      string `gorm:"unique;not null"`
	StoreID    uint32 `gorm:"not null"`
	Name       string `gorm:"not null"`
	Features   string `gorm:"not null"`
	ImgURL     string `gorm:"not null"`
	TimeToMake string `gorm:"not null"`
	CreatedAt  time.Time
}

type OrderItem struct {
	ID        uint32 `gorm:"primary_key;auto_increment"`
	OrderID   uint32 `gorm:"not null"`
	MenuID    uint32 `gorm:"not null"`
	Quantity  int    `gorm:"not null"`
	CreatedAt time.Time
}

type Device struct {
	ID           uint32 `gorm:"primary_key;auto_increment"`
	UserID       uint32 `gorm:"not null"`
	HashedData   []byte `gorm:"not null"`
	ScreenWidth  int    `gorm:"not null"`
	ScreenHeight int    `gorm:"not null"`
	IPAddress    string `gorm:"not null"`
	RemoteHost   string `gorm:"not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type StoreKeeperRequest struct {
	ID            uint32 `gorm:"primary_key;auto_increment"`
	StoreKeeperID uint32 `gorm:"not null"`
	RequestToken  string `gorm:"unique;not null"`
	CreatedAt     time.Time
	DeletedAt     *time.Time
}

type MenuDetail struct {
	MenuID      uint32 `gorm:"primary_key"`
	Remaining   int    `gorm:"not null"`
	TicketPrice int    `gorm:"not null"`
	Discount    int    `gorm:"not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Ticket struct {
	ID                uint32 `gorm:"primary_key"`
	YenPricePerTicket int    `gorm:"not null"`
}

type Barcode struct {
	ID          uint32 `gorm:"primary_key;auto_increment"`
	OrderID     uint32 `gorm:"not null"`
	BarcodeData string `gorm:"unique;not null"`
	CreatedAt   time.Time
}

type MenuAllergen struct {
	MenuID    uint32 `gorm:"primary_key"`
	Ebi       string `gorm:"not null"`
	Kani      string `gorm:"not null"`
	Komugi    string `gorm:"not null"`
	Kurumi    string `gorm:"not null"`
	Milk      string `gorm:"not null"`
	Peanut    string `gorm:"not null"`
	Soba      string `gorm:"not null"`
	Tamago    string `gorm:"not null"`
	CreatedAt time.Time
}
