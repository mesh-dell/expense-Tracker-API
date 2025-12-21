package database

import (
	"fmt"
	"log"

	"github.com/mesh-dell/expense-Tracker-API/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB(cfg config.Config) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser, cfg.DBPassword, cfg.DBAddress, cfg.DBName,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database")
	}
	return db
}
