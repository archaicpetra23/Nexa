package main

import (
	"log"
	"nexa/backend/internal/config"
	"nexa/backend/internal/database"
	"nexa/backend/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg.DBDSN)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	r := gin.Default()
	routes.Setup(r, db, cfg.JWTSecret)

	log.Printf("Server running on %s", cfg.Port)
	if err := r.Run(cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
