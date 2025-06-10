package server

import (
	"stock-management/internal/domain/usecases"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Server struct {
	db           *gorm.DB
	router       *gin.Engine
	authService  *usecases.AuthService
	stockService *usecases.StockService
}

func NewServer(db *gorm.DB, authService *usecases.AuthService, stockService *usecases.StockService) *Server {
	server := &Server{
		db:           db,
		router:       gin.Default(),
		authService:  authService,
		stockService: stockService,
	}

	server.setupCORS()
	server.setupRoutes()
	return server
}

func (s *Server) setupCORS() {
	s.router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:4200"}, // Angular default port
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
}

func (s *Server) setupRoutes() {
	// Auth routes
	auth := s.router.Group("/api/auth")
	{
		auth.POST("/login", s.handleLogin)
		auth.POST("/register", s.handleRegister)
		auth.POST("/logout", s.handleLogout)
	}

	// Stock routes
	stock := s.router.Group("/api/stock")
	{
		stock.POST("/import", s.handleImportStock)
		stock.POST("/export", s.handleExportStock)
		stock.GET("/current", s.handleGetCurrentStock)
		stock.GET("/movements", s.handleGetStockMovements)
		stock.GET("/summary", s.handleGetStockSummary)
	}
}

func (s *Server) Start() error {
	return s.router.Run(":8080")
}
