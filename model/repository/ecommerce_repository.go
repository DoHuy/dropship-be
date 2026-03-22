package repository

import (
	"dropshipbe/dropshipbe"
	model "dropshipbe/model/schema"

	"gorm.io/gorm"
)

type EcommerceRepository interface {
	GetProducts(request *dropshipbe.DefaultRequest) ([]model.Product, error)
}

type defaultEcommerceRepository struct {
	db *gorm.DB
}

func (d *defaultEcommerceRepository) GetProducts(request *dropshipbe.DefaultRequest) ([]model.Product, error) {
	var products []model.Product

	// Khởi tạo câu truy vấn từ model Product
	query := d.db.Model(&model.Product{})

	// 2. Preload các quan hệ cần thiết để tránh N+1 Query
	query = query.
		Preload("Country").
		Preload("Categories").
		Preload("Images").
		Preload("PriceTiers").
		Preload("Options.OptionValues").
		Preload("Variants.OptionValues")

	query = query.Where("status = ?", "active")

	if request.CountryCode != "" {
		query = query.Joins("JOIN countries ON countries.id = products.country_id").
			Where("countries.code = ?", request.CountryCode)
	}

	err := query.Find(&products).Error
	if err != nil {
		return nil, err // Trả về lỗi nếu có vấn đề với Database
	}

	return products, nil
}

func NewEcommerceRepository(db *gorm.DB) EcommerceRepository {
	return &defaultEcommerceRepository{db: db}
}
