package service

import (
	"context"

	"github.com/mesh-dell/expense-Tracker-API/internal/custom"
	"github.com/mesh-dell/expense-Tracker-API/internal/expenses"
	"github.com/mesh-dell/expense-Tracker-API/internal/expenses/dtos"
	"github.com/mesh-dell/expense-Tracker-API/internal/expenses/repository"
)

type ExpenseService struct {
	repo repository.IExpenseRepository
}

func NewExpenseService(r repository.IExpenseRepository) *ExpenseService {
	return &ExpenseService{
		repo: r,
	}
}

func (svc *ExpenseService) AddExpense(ctx context.Context, userID uint, req dtos.ExpenseRequestDto) (*expenses.Expense, error) {
	if !expenses.IsValidCategory(req.Category) {
		return nil, custom.ErrInvalidCategory
	}
	expense := &expenses.Expense{
		Title:    req.Title,
		Category: req.Category,
		UserID:   userID,
		Amount:   req.Amount,
		Date:     req.Date,
	}
	err := svc.repo.AddExpense(ctx, expense)
	if err != nil {
		return nil, err
	}
	return expense, err
}
func (svc *ExpenseService) GetAllExpensesForUser(ctx context.Context, userID uint, filter dtos.ExpenseFilter) ([]expenses.Expense, error) {
	exps, err := svc.repo.GetAllExpensesForUser(ctx, userID, filter)
	if err != nil {
		return nil, err
	}
	return exps, nil
}
func (svc *ExpenseService) GetExpenseByID(ctx context.Context, expenseID, userID uint) (expenses.Expense, error) {
	return svc.repo.GetExpenseByID(ctx, expenseID, userID)

}
func (svc *ExpenseService) RemoveExpense(ctx context.Context, expenseID, userID uint) error {
	return svc.repo.RemoveExpense(ctx, expenseID, userID)
}
func (svc *ExpenseService) UpdateExpense(ctx context.Context, req dtos.ExpenseRequestDto, expenseID, userID uint) (*expenses.Expense, error) {
	if !expenses.IsValidCategory(req.Category) {
		return nil, custom.ErrInvalidCategory
	}
	expense := &expenses.Expense{
		Title:    req.Title,
		Category: req.Category,
		UserID:   userID,
		Amount:   req.Amount,
		Date:     req.Date,
	}
	updatedExpense, err := svc.repo.UpdateExpense(ctx, *expense, expenseID, userID)
	if err != nil {
		return nil, err
	}
	return updatedExpense, err
}
