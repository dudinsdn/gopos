package models

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string `gorm:"size:100;not null"`
	Email    string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Role     string `gorm:"default:'cashier'"`
}

type Product struct {
	ID    uint    `gorm:"primaryKey"`
	Name  string  `gorm:"not null"`
	Price float64 `gorm:"not null"`
	Stock int     `gorm:"not null"`
}

type Transaction struct {
	ID       uint              `gorm:"primaryKey"`
	UserID   uint              `gorm:"not null"`
	User     User              `gorm:"foreignKey:UserID"`
	Total    float64           `gorm:"not null"`
	CreateAt int64             `gorm:"autoCreateTime:nano"`
	Items    []TransactionItem `gorm:"foreignKey:TransactionID"`
}

type TransactionItem struct {
	ID            uint    `gorm:"primaryKey"`
	TransactionID uint    `gorm:"not null"`
	ProductID     uint    `gorm:"not null"`
	Product       Product `gorm:"ProductID"`
	Quantity      int     `gorm:"not null"`
	Subtotal      float64 `gorm:"not null"`
}
