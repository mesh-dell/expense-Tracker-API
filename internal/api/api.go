package api

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mesh-dell/expense-Tracker-API/internal/api/middleware"
	"github.com/mesh-dell/expense-Tracker-API/internal/config"
	"github.com/mesh-dell/expense-Tracker-API/internal/database"
	"github.com/mesh-dell/expense-Tracker-API/internal/expenses"
	expensesHandler "github.com/mesh-dell/expense-Tracker-API/internal/expenses/handler"
	expensesRepository "github.com/mesh-dell/expense-Tracker-API/internal/expenses/repository"
	expensesService "github.com/mesh-dell/expense-Tracker-API/internal/expenses/service"
	"github.com/mesh-dell/expense-Tracker-API/internal/users"
	"github.com/mesh-dell/expense-Tracker-API/internal/users/handler"
	"github.com/mesh-dell/expense-Tracker-API/internal/users/repository"
	"github.com/mesh-dell/expense-Tracker-API/internal/users/service"
)

func InitServer(cfg config.Config) {
	// init db
	db := database.InitDB(cfg)
	db.AutoMigrate(&users.User{}, &expenses.Expense{}, &users.RefreshToken{})

	userRepo := repository.NewUserRepository(db)
	refreshRepo := repository.NewRefreshTokenRepository(db)
	refreshSvc := service.NewRefreshTokenService(refreshRepo)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService, cfg, refreshSvc)

	expensesRepo := expensesRepository.NewExpenseRepository(db)
	expensesSvc := expensesService.NewExpenseService(expensesRepo)
	expenseHandler := expensesHandler.NewExpenseHandler(expensesSvc)

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	protected := router.Group("/expenses")
	protected.Use(middleware.AuthMiddleware(cfg))
	{
		protected.POST("", expenseHandler.Create)
		protected.GET("/:id", expenseHandler.FindById)
		protected.GET("", expenseHandler.FindAllForUser)
		protected.PUT("/:id", expenseHandler.Update)
		protected.DELETE("/:id", expenseHandler.Delete)
	}
	// auth routes
	router.POST("/login", userHandler.Login)
	router.POST("/register", userHandler.Register)
	router.POST("/token/refresh", userHandler.RefreshToken)
	router.Run(":" + cfg.Port)
}
