package server

import (
	"stock-management/internal/domain/services"
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
	jwtService   services.JWTService
}

func NewServer(db *gorm.DB, authService *usecases.AuthService, stockService *usecases.StockService, jwtService services.JWTService) *Server {
	server := &Server{
		db:           db,
		router:       gin.Default(),
		authService:  authService,
		stockService: stockService,
		jwtService:   jwtService,
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

	// Product routes
	products := s.router.Group("/api/products")
	products.Use(AuthMiddleware(s.jwtService))
	{
		products.GET("", s.handleGetProducts)
		products.POST("", s.handleCreateProduct)
		products.PUT("/:id", s.handleUpdateProduct)
		products.DELETE("/:id", s.handleDeleteProduct)
	}

	// Category routes
	categories := s.router.Group("/api/categories")
	categories.Use(AuthMiddleware(s.jwtService))
	{
		categories.GET("", s.handleGetCategories)
		categories.POST("", s.handleCreateCategory)
		categories.PUT("/:id", s.handleUpdateCategory)
		categories.DELETE("/:id", s.handleDeleteCategory)
	}

	// Stock routes
	stock := s.router.Group("/api/stock")
	stock.Use(AuthMiddleware(s.jwtService))
	{
		stock.POST("/import", s.handleImportStock)
		stock.POST("/export", s.handleExportStock)
		stock.GET("/current", s.handleGetCurrentStock)
		stock.POST("/movements", s.handleGetStockMovements)
		stock.GET("/summary", s.handleGetStockSummary)
	}
}

func (s *Server) Start() error {
	return s.router.Run(":8080")
}
