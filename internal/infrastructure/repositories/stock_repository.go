package repositories

import (
	"context"
	"stock-management/internal/domain/models"
	"stock-management/internal/domain/repositories"
	"time"

	"gorm.io/gorm"
)

type stockRepository struct {
	db *gorm.DB
}

func NewStockRepository(db *gorm.DB) repositories.StockRepository {
	return &stockRepository{db: db}
}

func (r *stockRepository) CreateProduct(ctx context.Context, product *models.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

func (r *stockRepository) GetProduct(ctx context.Context, id uint) (*models.Product, error) {
	var product models.Product
	err := r.db.WithContext(ctx).First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *stockRepository) UpdateStock(ctx context.Context, stock *models.Stock) error {
	return r.db.WithContext(ctx).Save(stock).Error
}

func (r *stockRepository) GetStock(ctx context.Context, productID uint) (*models.Stock, error) {
	var stock models.Stock
	err := r.db.WithContext(ctx).Where("product_id = ?", productID).First(&stock).Error
	if err != nil {
		return nil, err
	}
	return &stock, nil
}

func (r *stockRepository) CreateMovement(ctx context.Context, movement *models.StockMovement) error {
	return r.db.WithContext(ctx).Create(movement).Error
}

func (r *stockRepository) GetMovements(ctx context.Context, startDate, endDate time.Time, productID *uint, category *string) ([]models.StockMovement, error) {
	query := r.db.WithContext(ctx).Model(&models.StockMovement{}).
		Where("date BETWEEN ? AND ?", startDate, endDate)

	if productID != nil {
		query = query.Where("product_id = ?", *productID)
	}

	if category != nil {
		query = query.Joins("JOIN products ON products.id = stock_movements.product_id").
			Where("products.category = ?", *category)
	}

	var movements []models.StockMovement
	err := query.Find(&movements).Error
	return movements, err
}

func (r *stockRepository) GetStockSummary(ctx context.Context) ([]models.Stock, error) {
	var stocks []models.Stock
	err := r.db.WithContext(ctx).Find(&stocks).Error
	return stocks, err
}
