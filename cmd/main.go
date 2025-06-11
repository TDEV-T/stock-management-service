package main

import (
	"log"
	"stock-management/config"
	"stock-management/internal/domain/repositories"
	"stock-management/internal/domain/services"
	"stock-management/internal/domain/usecases"
	"stock-management/internal/infrastructure/database"
	"stock-management/internal/infrastructure/server"
)

func main() {
	// Initialize database connection
	db, err := database.NewDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	//Load Config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	stockRepo := repositories.NewStockRepository(db)

	// Initialize services
	jwtService, err := services.NewJWTService(cfg.JWTSecret, "stock", 6400)
	authService := usecases.NewAuthService(userRepo, jwtService)
	stockService := usecases.NewStockService(stockRepo)

	// Initialize and start the server
	srv := server.NewServer(db, authService, stockService, jwtService)
	if err := srv.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
