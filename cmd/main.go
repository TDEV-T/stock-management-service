package main

import (
	"log"
	"stock-management/internal/domain/usecases"
	"stock-management/internal/infrastructure/database"
	"stock-management/internal/infrastructure/repositories"
	"stock-management/internal/infrastructure/server"
)

func main() {
	// Initialize database connection
	db, err := database.NewDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	stockRepo := repositories.NewStockRepository(db)

	// Initialize services
	authService := usecases.NewAuthService(userRepo)
	stockService := usecases.NewStockService(stockRepo)

	// Initialize and start the server
	srv := server.NewServer(db, authService, stockService)
	if err := srv.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
