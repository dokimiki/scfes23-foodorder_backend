package types

import (
	"time"
)

type AllergensList struct {
	Ebi    string `json:"ebi"`
	Kani   string `json:"kani"`
	Komugi string `json:"komugi"`
	Kurumi string `json:"kurumi"`
	Milk   string `json:"milk"`
	Peanut string `json:"peanut"`
	Soba   string `json:"soba"`
	Tamago string `json:"tamago"`
}

type Coupon struct {
	Kind string `json:"kind"`
}

type CouponItemIds struct {
	None         *string `json:"none"`
	Zero         *string `json:"0"`
	OneHundred   *string `json:"100"`
	TwoHundred   *string `json:"200"`
	ThreeHundred *string `json:"300"`
}

type MenuItem struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Price  int    `json:"price"`
	Image  string `json:"image"`
	IsShow bool   `json:"isShow"`
}

type CartItem struct {
	ID       string `json:"id"`
	Quantity int    `json:"quantity"`
}

type Order struct {
	ID            string     `json:"id"`
	IsMobileOrder bool       `json:"isMobileOrder"`
	NumberTag     int        `json:"numberTag"`
	Items         []CartItem `json:"items"`
}

type CompleteState struct {
	State string `json:"state"`
}

type CompleteInfo struct {
	Barcode      string     `json:"barcode"`
	CompleteTime time.Time  `json:"completeTime"`
	Items        []CartItem `json:"items"`
}

type OrderedPotato struct {
	ReceptionTime  time.Time
	CompletionTime time.Time
	Qty            int
	Order          struct {
		ID            string
		IsMobileOrder bool
		IsPaid        bool
		NumberTag     int
	}
}

type User struct {
	ID        string `json:"id"`
	IsOrdered bool   `json:"isOrdered"`
}

type InvitationStatus struct {
	IsInvited bool `json:"isInvited"`
}
