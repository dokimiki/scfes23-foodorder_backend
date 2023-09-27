package types

// AllergenContaminationStatus
type AllergenContaminationStatus string

const (
    AllergenContaminationStatusNotContains AllergenContaminationStatus = "NotContains"
    AllergenContaminationStatusContamination AllergenContaminationStatus = "Contamination"
    AllergenContaminationStatusContains AllergenContaminationStatus = "Contains"
)

// AllergensList
type AllergensList struct {
    Ebi AllergenContaminationStatus `json:"ebi"`
    Kani AllergenContaminationStatus `json:"kani"`
    Komugi AllergenContaminationStatus `json:"komugi"`
    Kurumi AllergenContaminationStatus `json:"kurumi"`
    Milk AllergenContaminationStatus `json:"milk"`
    Peanut AllergenContaminationStatus `json:"peanut"`
    Soba AllergenContaminationStatus `json:"soba"`
    Tamago AllergenContaminationStatus `json:"tamago"`
}

// MenuItem
type MenuItem struct {
    ID string `json:"id"`
    Name string `json:"name"`
    Price int `json:"price"`
    Image string `json:"image"`
}

// CartItem
type CartItem struct {
    ID string `json:"id"`
    Quantity int `json:"quantity"`
}

// Order
type Order struct {
    ID string `json:"id"`
    IsMobileOrder bool `json:"isMobileOrder"`
    NumberTag int `json:"numberTag"`
    Items []CartItem `json:"items"`
}

// OrderedPotato
type OrderedPotato struct {
    ReceptionTime time.Time
    CompletionTime time.Time
    Qty int
    Order struct {
        ID string
        IsMobileOrder bool
        NumberTag int
    }
}