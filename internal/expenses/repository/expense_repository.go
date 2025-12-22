package repository

import (
	"context"
	"errors"

	"github.com/mesh-dell/expense-Tracker-API/internal/custom"
	"github.com/mesh-dell/expense-Tracker-API/internal/expenses"
	"github.com/mesh-dell/expense-Tracker-API/internal/expenses/dtos"
	"gorm.io/gorm"
)

type ExpenseRepository struct {
	db *gorm.DB
}

// AddExpense implements [IExpenseRepository].
func (r *ExpenseRepository) AddExpense(ctx context.Context, expense *expenses.Expense) error {
	return gorm.G[expenses.Expense](r.db).Create(ctx, expense)
}

// GetAllExpensesForUser implements [IExpenseRepository].
func (r *ExpenseRepository) GetAllExpensesForUser(ctx context.Context, userID uint, filter dtos.ExpenseFilter) ([]expenses.Expense, error) {
	query := gorm.G[expenses.Expense](r.db).Where("user_id = ?", userID)
	if filter.StartDate != nil {
		query = query.Where("date >= ?", *filter.StartDate)
	}
	if filter.EndDate != nil {
		query = query.Where("date <= ?", *filter.EndDate)
	}
	return query.Find(ctx)
}

// GetExpenseByID implements [IExpenseRepository].
func (r *ExpenseRepository) GetExpenseByID(ctx context.Context, expenseID uint, userID uint) (expenses.Expense, error) {
	expense, err := gorm.G[expenses.Expense](r.db).Where("id = ? AND user_id = ?", expenseID, userID).First(ctx)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return expenses.Expense{}, custom.ErrExpenseNotFound
	}
	if err != nil {
		return expenses.Expense{}, err
	}
	return expense, nil
}

// RemoveExpense implements [IExpenseRepository].
func (r *ExpenseRepository) RemoveExpense(ctx context.Context, expenseID uint, userID uint) error {
	rowsAff, err := gorm.G[expenses.Expense](r.db).Where("id = ? AND user_id = ?", expenseID, userID).Delete(ctx)
	if err != nil {
		return err
	}
	if rowsAff == 0 {
		return custom.ErrExpenseNotFound
	}
	return nil
}

// UpdateExpense implements [IExpenseRepository].
func (r *ExpenseRepository) UpdateExpense(ctx context.Context, expense expenses.Expense, expenseID uint, userID uint) (*expenses.Expense, error) {
	rowsAff, err := gorm.G[expenses.Expense](r.db).Where("id = ? AND user_id = ?", expenseID, userID).Updates(ctx, expense)
	if err != nil {
		return nil, err
	}
	if rowsAff == 0 {
		return nil, custom.ErrExpenseNotFound
	}
	expenseUpdated, err := gorm.G[expenses.Expense](r.db).
		Where("id = ? AND user_id = ?", expenseID, userID).
		First(ctx)
	if err != nil {
		return nil, err
	}
	return &expenseUpdated, nil
}

func NewExpenseRepository(db *gorm.DB) IExpenseRepository {
	return &ExpenseRepository{
		db: db,
	}
}
