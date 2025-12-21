package users

import (
	"time"

	"github.com/mesh-dell/expense-Tracker-API/internal/expenses"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name         string `gorm:"not null"`
	Email        string `gorm:"not null;unique"`
	PasswordHash string `gorm:"not null"`
	Expenses     []expenses.Expense
}

type RefreshToken struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null;index"`
	JTI       string    `gorm:"uniqueIndex;size:500;not null"`
	ExpiresAt time.Time `gorm:"not null"`
}
