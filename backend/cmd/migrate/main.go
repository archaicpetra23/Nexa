package main

import (
	"fmt"
	"nexa/backend/internal/config"
	"nexa/backend/internal/database"
)

func main() {
	cfg := config.Load()
	_, err := database.Connect(cfg.DBDSN)
	if err != nil {
		fmt.Println("ERR:", err)
		return
	}
	fmt.Println("OK")
}
