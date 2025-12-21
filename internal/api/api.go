package api

import (
	"github.com/mesh-dell/expense-Tracker-API/internal/config"
	"github.com/mesh-dell/expense-Tracker-API/internal/database"
	"github.com/mesh-dell/expense-Tracker-API/internal/expenses"
	"github.com/mesh-dell/expense-Tracker-API/internal/users"
)

func InitServer(cfg config.Config) {
	// init db
	db := database.InitDB(cfg)
	db.AutoMigrate(&users.User{}, &expenses.Expense{})
}
