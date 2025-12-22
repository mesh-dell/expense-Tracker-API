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

type ExpenseCategory string

const (
	CategoryGroceries   ExpenseCategory = "Groceries"
	CategoryLeisure     ExpenseCategory = "Leisure"
	CategoryElectronics ExpenseCategory = "Electronics"
	CategoryUtilities   ExpenseCategory = "Utilities"
	CategoryClothing    ExpenseCategory = "Clothing"
	CategoryHealth      ExpenseCategory = "Health"
	CategoryOthers      ExpenseCategory = "Others"
)

func IsValidCategory(cat string) bool {
	switch ExpenseCategory(cat) {
	case CategoryGroceries,
		CategoryLeisure,
		CategoryElectronics,
		CategoryUtilities,
		CategoryClothing,
		CategoryHealth,
		CategoryOthers:
		return true
	default:
		return false
	}
}
