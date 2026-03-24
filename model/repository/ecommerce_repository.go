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
	GetFeaturedProducts(ctx context.Context, request *dropshipbe.DefaultRequest) ([]model.Product, error)
	GetNewProducts(ctx context.Context, request *dropshipbe.DefaultRequest) ([]model.Product, error)
	GetBannerItems(ctx context.Context, request *dropshipbe.DefaultRequest) ([]model.Banner, error)
	GetBlogBySlug(ctx context.Context, request *dropshipbe.GetBlogBySlugRequest) (*model.BlogPost, error)
	GetBlogItems(ctx context.Context, request *dropshipbe.DefaultRequest) ([]model.BlogPost, error)
	GetCategoryItems(ctx context.Context, request *dropshipbe.DefaultRequest) ([]model.Category, error)
}

type defaultEcommerceRepository struct {
	db    *gorm.DB
	cache cache.Cache // Nhận bộ nhớ đệm đã được tuỳ biến TTL từ bên ngoài
}

// GetCategoryItems implements [EcommerceRepository].
func (d *defaultEcommerceRepository) GetCategoryItems(ctx context.Context, request *dropshipbe.DefaultRequest) ([]model.Category, error) {
	var categories []model.Category
	cacheKey := constant.CategoryListKey(request.CountryCode)
	err := d.cache.TakeCtx(ctx, &categories, cacheKey, func(v any) error {
		query := d.db.Model(&model.Category{})
		if request.CountryCode != "" {
			query = query.Joins("JOIN countries ON countries.code = categories.country_code").
				Where("countries.code = ?", request.CountryCode)
		}
		return query.Find(v).Error
	})

	if err != nil {
		return nil, err
	}

	return categories, nil
}

// GetBlogItems implements [EcommerceRepository].
func (d *defaultEcommerceRepository) GetBlogItems(ctx context.Context, request *dropshipbe.DefaultRequest) ([]model.BlogPost, error) {
	var blogPosts []model.BlogPost
	cacheKey := constant.BlogPostListByCountryKey(request.CountryCode)

	err := d.cache.TakeCtx(ctx, &blogPosts, cacheKey, func(v any) error {
		query := d.db.Model(&model.BlogPost{}).
			Preload("Category").
			Preload("Country").
			Where("published_at IS NOT NULL")

		if request.CountryCode != "" {
			query = query.Joins("JOIN countries ON countries.code = blog_posts.country_code").
				Where("countries.code = ?", request.CountryCode)
		}

		return query.Find(v).Error
	})

	if err != nil {
		return nil, err
	}

	return blogPosts, nil
}

func (d *defaultEcommerceRepository) GetBlogBySlug(ctx context.Context, request *dropshipbe.GetBlogBySlugRequest) (*model.BlogPost, error) {
	var blogPost model.BlogPost
	cacheKey := constant.BlogPostBySlugKey(request.Slug, request.CountryCode)

	err := d.cache.TakeCtx(ctx, &blogPost, cacheKey, func(v any) error {
		query := d.db.Model(&model.BlogPost{})
		if request.CountryCode != "" {
			query = query.Joins("JOIN countries ON countries.code = blog_posts.country_code").
				Where("countries.code = ?", request.CountryCode)
		}

		return query.
			Preload("Category").
			Preload("Country").
			Where("slug = ?", request.Slug).First(v).Error
	})

	if err != nil {
		return nil, err
	}

	return &blogPost, nil
}

func (d *defaultEcommerceRepository) GetBannerItems(ctx context.Context, request *dropshipbe.DefaultRequest) ([]model.Banner, error) {
	var banners []model.Banner
	cacheKey := constant.BannerItemListKey()

	err := d.cache.TakeCtx(ctx, &banners, cacheKey, func(v any) error {
		return d.db.Model(&model.Banner{}).Where("is_active = ?", true).
			Find(v).Error
	})

	if err != nil {
		return nil, err
	}

	return banners, nil
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

func (d *defaultEcommerceRepository) GetFeaturedProducts(ctx context.Context, request *dropshipbe.DefaultRequest) ([]model.Product, error) {
	var products []model.Product

	cacheKey := constant.FeaturedProductListKey(request.CountryCode)

	err := d.cache.TakeCtx(ctx, &products, cacheKey, func(v any) error {

		query := d.db.Model(&model.Product{}).
			Preload("Country").
			Preload("Categories").
			Preload("Images").
			Preload("PriceTiers").
			Preload("Options.OptionValues").
			Preload("Variants.OptionValues").
			Where("status = ? AND is_featured = ?", "active", true)

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

func (d *defaultEcommerceRepository) GetNewProducts(ctx context.Context, request *dropshipbe.DefaultRequest) ([]model.Product, error) {
	var products []model.Product

	cacheKey := constant.NewProductListKey(request.CountryCode)

	err := d.cache.TakeCtx(ctx, &products, cacheKey, func(v any) error {

		query := d.db.Model(&model.Product{}).
			Preload("Country").
			Preload("Categories").
			Preload("Images").
			Preload("PriceTiers").
			Preload("Options.OptionValues").
			Preload("Variants.OptionValues").
			Where("status = ? AND is_new = ?", "active", true)

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
