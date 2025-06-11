package repositories

import (
	"context"
	"stock-management/internal/domain/models"
	"time"

	"gorm.io/gorm"
)

type StockRepository interface {
	CreateProduct(ctx context.Context, product *models.Product) error
	GetProduct(ctx context.Context, id uint) (*models.Product, error)
	GetProducts(ctx context.Context) ([]models.Product, error)
	UpdateProduct(ctx context.Context, product *models.Product) error
	DeleteProduct(ctx context.Context, id uint) error
	CreateCategory(ctx context.Context, category *models.Category) error
	GetCategory(ctx context.Context, id uint) (*models.Category, error)
	GetCategories(ctx context.Context) ([]models.Category, error)
	UpdateCategory(ctx context.Context, category *models.Category) error
	DeleteCategory(ctx context.Context, id uint) error

	GetStock(ctx context.Context, productID uint) (*models.Stock, error)
	CreateMovement(ctx context.Context, movement *models.StockMovement) error
	GetMovements(ctx context.Context, startDate, endDate time.Time, productID *uint, category *uint) ([]models.StockMovement, error)
	GetStockByProductID(productID uint) (*models.Stock, error)
	CreateStock(stock *models.Stock, movement *models.StockMovement) error
	UpdateStock(stock *models.Stock, movement *models.StockMovement) error
	GetStockMovements(startDate, endDate *time.Time, productID, categoryID *uint) ([]models.StockMovement, error)
	GetCurrentStock() ([]models.Stock, error)
	GetStockSummary() ([]models.Stock, error)
}

type stockRepository struct {
	db *gorm.DB
}

func NewStockRepository(db *gorm.DB) StockRepository {
	return &stockRepository{db: db}
}

func (r *stockRepository) CreateProduct(ctx context.Context, product *models.Product) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create the product
		if err := tx.Create(product).Error; err != nil {
			return err
		}

		// Create initial stock record
		stock := &models.Stock{
			ProductID: product.ID,
			Quantity:  0,
		}
		return tx.Create(stock).Error
	})
}

func (r *stockRepository) GetProduct(ctx context.Context, id uint) (*models.Product, error) {
	var product models.Product
	err := r.db.WithContext(ctx).First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *stockRepository) GetProducts(ctx context.Context) ([]models.Product, error) {
	var products []models.Product
	err := r.db.WithContext(ctx).Preload("Category").Find(&products).Error
	return products, err
}

func (r *stockRepository) UpdateProduct(ctx context.Context, product *models.Product) error {
	return r.db.WithContext(ctx).Save(product).Error
}

func (r *stockRepository) DeleteProduct(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Delete associated stock records
		if err := tx.Where("product_id = ?", id).Delete(&models.Stock{}).Error; err != nil {
			return err
		}

		// Delete associated stock movements
		if err := tx.Where("product_id = ?", id).Delete(&models.StockMovement{}).Error; err != nil {
			return err
		}

		// Delete the product
		return tx.Delete(&models.Product{}, id).Error
	})
}

func (r *stockRepository) CreateCategory(ctx context.Context, category *models.Category) error {
	return r.db.WithContext(ctx).Create(category).Error
}

func (r *stockRepository) GetCategory(ctx context.Context, id uint) (*models.Category, error) {
	var category models.Category
	err := r.db.WithContext(ctx).First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *stockRepository) GetCategories(ctx context.Context) ([]models.Category, error) {
	var categories []models.Category
	err := r.db.WithContext(ctx).Find(&categories).Error
	return categories, err
}

func (r *stockRepository) UpdateCategory(ctx context.Context, category *models.Category) error {
	return r.db.WithContext(ctx).Save(category).Error
}

func (r *stockRepository) DeleteCategory(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Update products to remove category reference
		if err := tx.Model(&models.Product{}).Where("category_id = ?", id).Update("category_id", nil).Error; err != nil {
			return err
		}

		// Delete the category
		return tx.Delete(&models.Category{}, id).Error
	})
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

func (r *stockRepository) GetMovements(ctx context.Context, startDate, endDate time.Time, productID *uint, category *uint) ([]models.StockMovement, error) {
	query := r.db.WithContext(ctx).Model(&models.StockMovement{}).
		Preload("Product").
		Preload("Product.Category").
		Preload("User")

	// Only add date filter if both dates are provided
	if !startDate.IsZero() && !endDate.IsZero() {
		query = query.Where("date BETWEEN ? AND ?", startDate, endDate)
	}

	if productID != nil {
		query = query.Where("product_id = ?", *productID)
	}

	if category != nil {
		query = query.Joins("JOIN products ON products.id = stock_movements.product_id").
			Where("products.category_id = ?", *category)
	}

	var movements []models.StockMovement
	err := query.Find(&movements).Error
	return movements, err
}

func (r *stockRepository) GetStockByProductID(productID uint) (*models.Stock, error) {
	var stock models.Stock
	err := r.db.Where("product_id = ?", productID).First(&stock).Error
	if err != nil {
		return nil, err
	}
	return &stock, nil
}

func (r *stockRepository) CreateStock(stock *models.Stock, movement *models.StockMovement) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Create stock record
		if err := tx.Create(stock).Error; err != nil {
			return err
		}

		// Create movement record
		if err := tx.Create(movement).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *stockRepository) UpdateStock(stock *models.Stock, movement *models.StockMovement) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Update stock record
		if err := tx.Save(stock).Error; err != nil {
			return err
		}

		// Create movement record
		if err := tx.Create(movement).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *stockRepository) GetStockMovements(startDate, endDate *time.Time, productID, categoryID *uint) ([]models.StockMovement, error) {
	query := r.db.Model(&models.StockMovement{}).
		Preload("Product").
		Preload("Product.Category").
		Preload("User")

	if startDate != nil {
		query = query.Where("date >= ?", startDate)
	}
	if endDate != nil {
		query = query.Where("date <= ?", endDate)
	}
	if productID != nil {
		query = query.Where("product_id = ?", productID)
	}
	if categoryID != nil {
		query = query.Joins("JOIN products ON stock_movements.product_id = products.id").
			Where("products.category_id = ?", categoryID)
	}

	var movements []models.StockMovement
	err := query.Order("date DESC").Find(&movements).Error
	return movements, err
}

func (r *stockRepository) GetCurrentStock() ([]models.Stock, error) {
	var stocks []models.Stock
	err := r.db.Preload("Product").
		Preload("Product.Category").
		Find(&stocks).Error
	return stocks, err
}

func (r *stockRepository) GetStockSummary() ([]models.Stock, error) {
	var stocks []models.Stock
	err := r.db.Preload("Product").
		Preload("Product.Category").
		Find(&stocks).Error
	return stocks, err
}
