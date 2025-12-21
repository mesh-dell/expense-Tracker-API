package expenses

import (
	"time"

	"gorm.io/gorm"
)

type Expense struct {
	gorm.Model
	Title    string    `gorm:"not null"`
	Category string    `gorm:"not null"`
	UserID   uint      `gorm:"not null"`
	Amount   float64   `gorm:"not null"`
	Date     time.Time `gorm:"not null"`
}
