package main

import (
	"log"

	"github.com/mesh-dell/expense-Tracker-API/internal/api"
	"github.com/mesh-dell/expense-Tracker-API/internal/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("error load config: %v", err)
	}
	api.InitServer(cfg)
}
