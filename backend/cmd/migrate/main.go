package main

import (
	"fmt"
	"nexa/backend/internal/database"
	"nexa/backend/internal/config"
)

func main() {
	cfg := config.Load()
	_, err := database.Connect(cfg.DBDSN)
	if err != nil { fmt.Println("ERR:", err); return }
	fmt.Println("OK")
}
