package usecases

import (
	"context"
	"errors"
	"stock-management/internal/domain/models"
	"stock-management/internal/domain/repositories"
	"time"
)

type StockService struct {
	stockRepo repositories.StockRepository
}

func NewStockService(stockRepo repositories.StockRepository) *StockService {
	return &StockService{stockRepo: stockRepo}
}

func (s *StockService) CreateProduct(ctx context.Context, product *models.Product) error {
	if product.Name == "" {
		return errors.New("product name is required")
	}
	if product.CategoryID == 0 {
		return errors.New("category is required")
	}
	return s.stockRepo.CreateProduct(ctx, product)
}

func (s *StockService) UpdateProduct(ctx context.Context, product *models.Product) error {
	if product.ID == 0 {
		return errors.New("product ID is required")
	}
	if product.Name == "" {
		return errors.New("product name is required")
	}
	if product.CategoryID == 0 {
		return errors.New("category is required")
	}

	// Verify product exists
	existingProduct, err := s.stockRepo.GetProduct(ctx, product.ID)
	if err != nil {
		return err
	}
	if existingProduct == nil {
		return errors.New("product not found")
	}

	return s.stockRepo.UpdateProduct(ctx, product)
}

func (s *StockService) GetStockByProductID(ctx context.Context, productID uint) (*models.Stock, error) {
	return s.stockRepo.GetStock(ctx, productID)
}

func (s *StockService) DeleteProduct(ctx context.Context, id uint) error {
	// Verify product exists
	existingProduct, err := s.stockRepo.GetProduct(ctx, id)
	if err != nil {
		return err
	}
	if existingProduct == nil {
		return errors.New("product not found")
	}

	return s.stockRepo.DeleteProduct(ctx, id)
}

func (s *StockService) ImportStock(ctx context.Context, productID uint, quantity int, userID uint, notes string) error {
	if quantity <= 0 {
		return errors.New("quantity must be greater than 0")
	}

	// Create stock movement record
	movement := &models.StockMovement{
		ProductID: productID,
		UserID:    userID,
		Type:      "import",
		Quantity:  quantity,
		Date:      time.Now(),
		Notes:     notes,
	}

	// Update stock quantity
	stock, err := s.stockRepo.GetStock(ctx, productID)
	if err != nil {
		// If stock doesn't exist, create new stock
		stock = &models.Stock{
			ProductID: productID,
			Quantity:  quantity,
		}
		return s.stockRepo.CreateStock(stock, movement)
	}

	// Update existing stock
	stock.Quantity += quantity
	return s.stockRepo.UpdateStock(stock, movement)
}

func (s *StockService) ExportStock(ctx context.Context, productID uint, quantity int, userID uint, notes string) error {
	if quantity <= 0 {
		return errors.New("quantity must be greater than 0")
	}

	// Check if we have enough stock
	stock, err := s.stockRepo.GetStock(ctx, productID)
	if err != nil {
		return errors.New("product not found in stock")
	}

	if stock.Quantity < quantity {
		return errors.New("insufficient stock")
	}

	// Create stock movement record
	movement := &models.StockMovement{
		ProductID: productID,
		UserID:    userID,
		Type:      "export",
		Quantity:  quantity,
		Date:      time.Now(),
		Notes:     notes,
	}

	// Update stock quantity
	stock.Quantity -= quantity
	return s.stockRepo.UpdateStock(stock, movement)
}

func (s *StockService) GetStockMovements(ctx context.Context, startDate, endDate time.Time, productID *uint, category *uint) ([]models.MovementDTO, error) {
	movements, err := s.stockRepo.GetMovements(ctx, startDate, endDate, productID, category)
	if err != nil {
		return nil, err
	}

	var dtos []models.MovementDTO
	for _, movement := range movements {
		dto := models.MovementDTO{
			Type:     movement.Type,
			Quantity: movement.Quantity,
			Date:     movement.Date,
			Notes:    movement.Notes,
		}

		// Map product info
		if movement.Product.ID != 0 {
			dto.Product.Name = movement.Product.Name
			dto.Product.ImageURL = movement.Product.ImageURL
		}

		// Map user info
		if movement.User.ID != 0 {
			dto.User.Username = movement.User.Username
		}

		dtos = append(dtos, dto)
	}

	return dtos, nil
}

func (s *StockService) GetStockSummary(ctx context.Context) ([]models.Stock, error) {
	return s.stockRepo.GetStockSummary()
}

func (s *StockService) GetCategories(ctx context.Context) ([]models.Category, error) {
	return s.stockRepo.GetCategories(ctx)
}

func (s *StockService) CreateCategory(ctx context.Context, category *models.Category) error {
	if category.Name == "" {
		return errors.New("category name is required")
	}
	return s.stockRepo.CreateCategory(ctx, category)
}

func (s *StockService) UpdateCategory(ctx context.Context, category *models.Category) error {
	if category.ID == 0 {
		return errors.New("category ID is required")
	}
	if category.Name == "" {
		return errors.New("category name is required")
	}

	// Verify category exists
	existingCategory, err := s.stockRepo.GetCategory(ctx, category.ID)
	if err != nil {
		return err
	}
	if existingCategory == nil {
		return errors.New("category not found")
	}

	return s.stockRepo.UpdateCategory(ctx, category)
}

func (s *StockService) DeleteCategory(ctx context.Context, id uint) error {
	// Verify category exists
	existingCategory, err := s.stockRepo.GetCategory(ctx, id)
	if err != nil {
		return err
	}
	if existingCategory == nil {
		return errors.New("category not found")
	}

	return s.stockRepo.DeleteCategory(ctx, id)
}

func (s *StockService) GetProducts(ctx context.Context) ([]models.ProductDTO, error) {
	products, err := s.stockRepo.GetProducts(ctx)
	if err != nil {
		return nil, err
	}

	var productDTOs []models.ProductDTO
	for _, product := range products {
		stock, _ := s.stockRepo.GetStock(ctx, product.ID)
		quantity := 0
		if stock != nil {
			quantity = stock.Quantity
		}

		productDTO := models.ProductDTO{
			ID:          product.ID,
			Name:        product.Name,
			ImageURL:    product.ImageURL,
			Description: product.Description,
			CategoryID:  product.CategoryID,
			Category:    product.Category,
			SKU:         product.SKU,
			CreatedAt:   product.CreatedAt,
			UpdatedAt:   product.UpdatedAt,
			Quantity:    quantity,
		}
		productDTOs = append(productDTOs, productDTO)
	}

	return productDTOs, nil
}
