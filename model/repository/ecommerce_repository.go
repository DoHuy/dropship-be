package repository

import (
	"context"

	"dropshipbe/common/constant"
	"dropshipbe/dropshipbe"
	model "dropshipbe/model/schema"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"gorm.io/gorm"
)

type EcommerceRepository interface {
	GetProducts(ctx context.Context, request *dropshipbe.DefaultRequest) ([]model.Product, error)
}

type defaultEcommerceRepository struct {
	db    *gorm.DB
	cache cache.Cache // Nhận bộ nhớ đệm đã được tuỳ biến TTL từ bên ngoài
}

func NewEcommerceRepository(db *gorm.DB, c cache.Cache) EcommerceRepository {
	return &defaultEcommerceRepository{
		db:    db,
		cache: c,
	}
}

func (d *defaultEcommerceRepository) GetProducts(ctx context.Context, request *dropshipbe.DefaultRequest) ([]model.Product, error) {
	var products []model.Product

	cacheKey := constant.ProductListByCountryKey(request.CountryCode)

	err := d.cache.TakeCtx(ctx, &products, cacheKey, func(v any) error {

		query := d.db.Model(&model.Product{}).
			Preload("Country").
			Preload("Categories").
			Preload("Images").
			Preload("PriceTiers").
			Preload("Options.OptionValues").
			Preload("Variants.OptionValues").
			Where("status = ?", "active")

		if request.CountryCode != "" {
			query = query.Joins("JOIN countries ON countries.code = products.country_code").
				Where("countries.code = ?", request.CountryCode)
		}

		return query.Find(v).Error
	})

	if err != nil {
		return nil, err
	}

	return products, nil
}
