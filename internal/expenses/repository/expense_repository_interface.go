package repository

import (
	"context"

	"github.com/mesh-dell/expense-Tracker-API/internal/expenses"
	"github.com/mesh-dell/expense-Tracker-API/internal/expenses/dtos"
)

type IExpenseRepository interface {
	AddExpense(ctx context.Context, expense *expenses.Expense) error
	GetExpenseByID(ctx context.Context, expenseID uint, userID uint) (expenses.Expense, error)
	GetAllExpensesForUser(ctx context.Context, userID uint, filter dtos.ExpenseFilter) ([]expenses.Expense, error)
	UpdateExpense(ctx context.Context, expense expenses.Expense, expenseID uint, userID uint) (*expenses.Expense, error)
	RemoveExpense(ctx context.Context, expenseID uint, userID uint) error
}
