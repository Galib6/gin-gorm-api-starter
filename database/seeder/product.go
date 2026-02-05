package seeder

import (
	"myapp/core/entity"
	"myapp/support/logger"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

func ProductSeeder(db *gorm.DB) error {
	// Create categories first
	categories := []entity.Category{
		{
			ID:          uuid.MustParse("a0000000-0000-0000-0000-000000000001"),
			Name:        "Electronics",
			Description: "Electronic devices and gadgets",
		},
		{
			ID:          uuid.MustParse("a0000000-0000-0000-0000-000000000002"),
			Name:        "Clothing",
			Description: "Apparel and fashion items",
		},
		{
			ID:          uuid.MustParse("a0000000-0000-0000-0000-000000000003"),
			Name:        "Books",
			Description: "Books and literature",
		},
	}

	for _, category := range categories {
		var existing entity.Category
		if err := db.Where("id = ?", category.ID).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&category).Error; err != nil {
					logger.Error("Error seeding category: %v", err)
					return err
				}
				logger.Debug("Category seeded: %s", category.Name)
			}
		}
	}

	// Create products
	electronicsID := uuid.MustParse("a0000000-0000-0000-0000-000000000001")
	clothingID := uuid.MustParse("a0000000-0000-0000-0000-000000000002")
	booksID := uuid.MustParse("a0000000-0000-0000-0000-000000000003")

	products := []entity.Product{
		{
			ID:          uuid.MustParse("b0000000-0000-0000-0000-000000000001"),
			Name:        "Smartphone Pro",
			Description: "Latest flagship smartphone with advanced features",
			SKU:         "ELEC-PHONE-001",
			Price:       decimal.NewFromFloat(999.99),
			Stock:       50,
			CategoryID:  &electronicsID,
			IsActive:    true,
		},
		{
			ID:          uuid.MustParse("b0000000-0000-0000-0000-000000000002"),
			Name:        "Laptop Ultra",
			Description: "High-performance laptop for professionals",
			SKU:         "ELEC-LAPTOP-001",
			Price:       decimal.NewFromFloat(1499.99),
			Stock:       30,
			CategoryID:  &electronicsID,
			IsActive:    true,
		},
		{
			ID:          uuid.MustParse("b0000000-0000-0000-0000-000000000003"),
			Name:        "Cotton T-Shirt",
			Description: "Comfortable cotton t-shirt",
			SKU:         "CLOTH-TSHIRT-001",
			Price:       decimal.NewFromFloat(29.99),
			Stock:       200,
			CategoryID:  &clothingID,
			IsActive:    true,
		},
		{
			ID:          uuid.MustParse("b0000000-0000-0000-0000-000000000004"),
			Name:        "Programming Guide",
			Description: "Comprehensive guide to programming",
			SKU:         "BOOK-PROG-001",
			Price:       decimal.NewFromFloat(49.99),
			Stock:       5, // Low stock for testing
			CategoryID:  &booksID,
			IsActive:    true,
		},
		{
			ID:          uuid.MustParse("b0000000-0000-0000-0000-000000000005"),
			Name:        "Wireless Headphones",
			Description: "Premium wireless headphones with noise cancellation",
			SKU:         "ELEC-HEADPHONE-001",
			Price:       decimal.NewFromFloat(299.99),
			Stock:       3, // Low stock for testing
			CategoryID:  &electronicsID,
			IsActive:    false, // Inactive product for testing
		},
	}

	for _, product := range products {
		var existing entity.Product
		if err := db.Where("id = ?", product.ID).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&product).Error; err != nil {
					logger.Error("Error seeding product: %v", err)
					return err
				}
				logger.Debug("Product seeded: %s", product.Name)
			}
		}
	}

	return nil
}
