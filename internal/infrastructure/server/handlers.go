package server

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Auth handlers
func (s *Server) handleLogin(c *gin.Context) {
	var loginReq struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := s.authService.Login(c.Request.Context(), loginReq.Username, loginReq.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user":    user,
	})
}

func (s *Server) handleRegister(c *gin.Context) {
	var registerReq struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&registerReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := s.authService.Register(c.Request.Context(), registerReq.Username, registerReq.Password, registerReq.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Registration successful"})
}

func (s *Server) handleLogout(c *gin.Context) {
	// TODO: Implement token invalidation if using JWT
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

// Stock handlers
func (s *Server) handleImportStock(c *gin.Context) {
	var importReq struct {
		ProductID uint   `json:"product_id" binding:"required"`
		Quantity  int    `json:"quantity" binding:"required,min=1"`
		Notes     string `json:"notes"`
	}

	if err := c.ShouldBindJSON(&importReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context (assuming it's set by auth middleware)
	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	err := s.stockService.ImportStock(c.Request.Context(), importReq.ProductID, importReq.Quantity, userID, importReq.Notes)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Stock imported successfully"})
}

func (s *Server) handleExportStock(c *gin.Context) {
	var exportReq struct {
		ProductID uint   `json:"product_id" binding:"required"`
		Quantity  int    `json:"quantity" binding:"required,min=1"`
		Notes     string `json:"notes"`
	}

	if err := c.ShouldBindJSON(&exportReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	err := s.stockService.ExportStock(c.Request.Context(), exportReq.ProductID, exportReq.Quantity, userID, exportReq.Notes)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Stock exported successfully"})
}

func (s *Server) handleGetCurrentStock(c *gin.Context) {
	stock, err := s.stockService.GetStockSummary(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stock)
}

func (s *Server) handleGetStockMovements(c *gin.Context) {
	startDate, _ := time.Parse("2006-01-02", c.Query("start_date"))
	endDate, _ := time.Parse("2006-01-02", c.Query("end_date"))

	var productID *uint
	if pid := c.Query("product_id"); pid != "" {
		id := uint(0)
		productID = &id
	}

	var category *string
	if cat := c.Query("category"); cat != "" {
		category = &cat
	}

	movements, err := s.stockService.GetStockMovements(c.Request.Context(), startDate, endDate, productID, category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, movements)
}

func (s *Server) handleGetStockSummary(c *gin.Context) {
	summary, err := s.stockService.GetStockSummary(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}
