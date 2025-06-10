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

func (s *StockService) ImportStock(ctx context.Context, productID uint, quantity int, userID uint, notes string) error {
	stock, err := s.stockRepo.GetStock(ctx, productID)
	if err != nil {
		return err
	}

	movement := &models.StockMovement{
		ProductID: productID,
		UserID:    userID,
		Type:      "import",
		Quantity:  quantity,
		Date:      time.Now(),
		Notes:     notes,
	}

	if err := s.stockRepo.CreateMovement(ctx, movement); err != nil {
		return err
	}

	stock.Quantity += quantity
	return s.stockRepo.UpdateStock(ctx, stock)
}

func (s *StockService) ExportStock(ctx context.Context, productID uint, quantity int, userID uint, notes string) error {
	stock, err := s.stockRepo.GetStock(ctx, productID)
	if err != nil {
		return err
	}

	if stock.Quantity < quantity {
		return errors.New("insufficient stock")
	}

	movement := &models.StockMovement{
		ProductID: productID,
		UserID:    userID,
		Type:      "export",
		Quantity:  quantity,
		Date:      time.Now(),
		Notes:     notes,
	}

	if err := s.stockRepo.CreateMovement(ctx, movement); err != nil {
		return err
	}

	stock.Quantity -= quantity
	return s.stockRepo.UpdateStock(ctx, stock)
}

func (s *StockService) GetStockMovements(ctx context.Context, startDate, endDate time.Time, productID *uint, category *string) ([]models.StockMovement, error) {
	return s.stockRepo.GetMovements(ctx, startDate, endDate, productID, category)
}

func (s *StockService) GetStockSummary(ctx context.Context) ([]models.Stock, error) {
	return s.stockRepo.GetStockSummary(ctx)
}
