package users

import (
	"github.com/mesh-dell/expense-Tracker-API/internal/expenses"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name         string `gorm:"not null"`
	Email        string `gorm:"not null"`
	PasswordHash string `gorm:"not null"`
	Expenses     []expenses.Expense
}
