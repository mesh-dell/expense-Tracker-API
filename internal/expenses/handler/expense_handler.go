package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mesh-dell/expense-Tracker-API/internal/custom"
	"github.com/mesh-dell/expense-Tracker-API/internal/expenses/dtos"
	"github.com/mesh-dell/expense-Tracker-API/internal/expenses/service"
)

type ExpenseHandler struct {
	svc *service.ExpenseService
}

func NewExpenseHandler(svc *service.ExpenseService) *ExpenseHandler {
	return &ExpenseHandler{
		svc: svc,
	}
}

func (h *ExpenseHandler) Create(c *gin.Context) {
	userIDstr, exists := c.Get("userID")
	if !exists {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := userIDstr.(uint)

	var req dtos.ExpenseRequestDto
	if err := c.ShouldBindJSON(&req); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": fmt.Sprint("invalid request body: ", err.Error())})
		return
	}
	expense, err := h.svc.AddExpense(c.Request.Context(), userID, req)
	if errors.Is(err, custom.ErrInvalidCategory) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid expense category"})
		return
	}
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusCreated, dtos.ExpenseResponseDto{
		ID:       expense.ID,
		Title:    expense.Title,
		Category: expense.Category,
		Amount:   expense.Amount,
		Date:     expense.Date,
	})
}
func (h *ExpenseHandler) FindById(c *gin.Context) {
	userIDstr, exists := c.Get("userID")
	if !exists {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := userIDstr.(uint)

	expenseIDStr := c.Param("id")
	expenseID, err := strconv.ParseUint(expenseIDStr, 10, 64)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	expense, err := h.svc.GetExpenseByID(c.Request.Context(), uint(expenseID), userID)
	if errors.Is(err, custom.ErrExpenseNotFound) {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, dtos.ExpenseResponseDto{
		ID:       expense.ID,
		Title:    expense.Title,
		Category: expense.Category,
		Amount:   expense.Amount,
		Date:     expense.Date,
	})
}
func (h *ExpenseHandler) FindAllForUser(c *gin.Context) {
	userIDstr, exists := c.Get("userID")
	if !exists {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := userIDstr.(uint)

	filterType := c.Query("filter")

	var filter dtos.ExpenseFilter

	switch filterType {
	case "week":
		filter = PastWeek()
	case "month":
		filter = PastMonth()
	case "3months":
		filter = PastThreeMonths()
	case "custom":
		start, _ := time.Parse("2006-01-02", c.Query("start"))
		end, _ := time.Parse("2006-01-02", c.Query("end"))
		filter = CustomRange(start, end)
	}

	exps, err := h.svc.GetAllExpensesForUser(c.Request.Context(), userID, filter)
	var expsRes []dtos.ExpenseResponseDto

	for _, expense := range exps {
		expsRes = append(expsRes, dtos.ExpenseResponseDto{
			ID:       expense.ID,
			Title:    expense.Title,
			Category: expense.Category,
			Amount:   expense.Amount,
			Date:     expense.Date,
		})
	}
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"expenses": expsRes})
}
func (h *ExpenseHandler) Delete(c *gin.Context) {
	userIDstr, exists := c.Get("userID")
	if !exists {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := userIDstr.(uint)

	expenseIDStr := c.Param("id")
	expenseID, err := strconv.ParseUint(expenseIDStr, 10, 64)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	err = h.svc.RemoveExpense(c.Request.Context(), uint(expenseID), userID)
	if errors.Is(err, custom.ErrExpenseNotFound) {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusNoContent, gin.H{"message": "expense deleted successfully"})
}
func (h *ExpenseHandler) Update(c *gin.Context) {
	userIDstr, exists := c.Get("userID")
	if !exists {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := userIDstr.(uint)

	expenseIDStr := c.Param("id")
	expenseID, err := strconv.ParseUint(expenseIDStr, 10, 64)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req dtos.ExpenseRequestDto
	if err := c.ShouldBindJSON(&req); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": fmt.Sprint("invalid request body: ", err.Error())})
		return
	}
	expense, err := h.svc.UpdateExpense(c.Request.Context(), req, uint(expenseID), userID)
	if errors.Is(err, custom.ErrInvalidCategory) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid expense category"})
		return
	}
	if errors.Is(err, custom.ErrExpenseNotFound) {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(*expense)
	c.IndentedJSON(http.StatusOK, dtos.ExpenseResponseDto{
		ID:       expense.ID,
		Title:    expense.Title,
		Category: expense.Category,
		Amount:   expense.Amount,
		Date:     expense.Date,
	})
}

func PastWeek() dtos.ExpenseFilter {
	start := time.Now().AddDate(0, 0, -7)
	return dtos.ExpenseFilter{StartDate: &start}
}

func PastMonth() dtos.ExpenseFilter {
	start := time.Now().AddDate(0, -1, 0)
	return dtos.ExpenseFilter{StartDate: &start}
}

func PastThreeMonths() dtos.ExpenseFilter {
	start := time.Now().AddDate(0, -3, 0)
	return dtos.ExpenseFilter{StartDate: &start}
}

func CustomRange(start, end time.Time) dtos.ExpenseFilter {
	return dtos.ExpenseFilter{
		StartDate: &start,
		EndDate:   &end,
	}
}
