package dtos

import (
	"time"
)

type ExpenseRequestDto struct {
	Title    string    `json:"title" binding:"required"`
	Category string    `json:"category" binding:"required"`
	Amount   float64   `json:"amount" binding:"required,gt=0"`
	Date     time.Time `json:"date" binding:"required"`
}
type ExpenseResponseDto struct {
	ID       uint      `json:"id"`
	Title    string    `json:"title"`
	Category string    `json:"category"`
	Amount   float64   `json:"amount"`
	Date     time.Time `json:"date"`
}
type ExpenseFilter struct {
	StartDate *time.Time
	EndDate   *time.Time
}
