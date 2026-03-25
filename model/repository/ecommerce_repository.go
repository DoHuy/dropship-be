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
	GetShop(ctx context.Context, request *dropshipbe.ShopSearchParams) ([]model.Product, error)
	GetFeaturedProducts(ctx context.Context, request *dropshipbe.DefaultRequest) ([]model.Product, error)
	GetNewProducts(ctx context.Context, request *dropshipbe.DefaultRequest) ([]model.Product, error)
	GetRelatedProducts(ctx context.Context, request *dropshipbe.GetRelatedProductsRequest) ([]model.Product, error)
	GetBannerItems(ctx context.Context, request *dropshipbe.DefaultRequest) ([]model.Banner, error)
	GetBlogBySlug(ctx context.Context, request *dropshipbe.GetBlogBySlugRequest) (*model.BlogPost, error)
	GetBlogItems(ctx context.Context, request *dropshipbe.DefaultRequest) ([]model.BlogPost, error)
	GetCategoryItems(ctx context.Context, request *dropshipbe.DefaultRequest) ([]model.Category, error)
	GetProductBySlug(ctx context.Context, request *dropshipbe.GetProductBySlugRequest) (*model.Product, error)
	GetProductFaqs(ctx context.Context, request *dropshipbe.GetProductFaqsRequest) ([]model.ProductFAQ, error)
	GetProductReviews(ctx context.Context, request *dropshipbe.GetProductReviewsRequest) ([]model.ProductReview, error)
	GetProductsByCategory(ctx context.Context, request *dropshipbe.GetProductsByCategoryRequest) ([]model.Product, error)
	GetSliderItems(ctx context.Context, request *dropshipbe.DefaultRequest) ([]model.Slider, error)
	GetSocialProductVideos(ctx context.Context, request *dropshipbe.GetSocialProductVideoRequest) ([]model.ProductImage, error)
	GetVideoBanner(ctx context.Context, request *dropshipbe.DefaultRequest) (*model.Banner, error)
	CreateProductReview(ctx context.Context, request *dropshipbe.CreateProductReviewRequest) (*model.ProductReview, error)
}

type defaultEcommerceRepository struct {
	db    *gorm.DB
	cache cache.Cache // Nhận bộ nhớ đệm đã được tuỳ biến TTL từ bên ngoài
}

// CreateProductReview implements [EcommerceRepository].
func (d *defaultEcommerceRepository) CreateProductReview(ctx context.Context, request *dropshipbe.CreateProductReviewRequest) (*model.ProductReview, error) {
	review := &model.ProductReview{
		ProductID:    request.ProductId,
		AuthorName:   request.Name,
		AuthorEmail:  request.Email,
		AuthorAvatar: request.Avatar,
		Rating:       int(request.Rating),
		Content:      request.Comment,
		IsVerified:   true,
		Media: &model.ReviewMedia{
			Images: request.Images,
			Videos: request.Videos,
		},
		Status: "active",
	}

	if err := d.db.Create(review).Error; err != nil {
		return nil, err
	}

	return review, nil
}

// GetSocialProductVideos implements [EcommerceRepository].
func (d *defaultEcommerceRepository) GetSocialProductVideos(ctx context.Context, request *dropshipbe.GetSocialProductVideoRequest) ([]model.ProductImage, error) {
	var videos []model.ProductImage
	cacheKey := constant.SocialProductVideoListKey(request.Id, request.CountryCode)

	err := d.cache.TakeCtx(ctx, &videos, cacheKey, func(v any) error {
		query := d.db.Model(&model.ProductImage{}).Where("product_id = ? AND media_type = ?", request.Id, "social_video")
		return query.Find(v).Error
	})

	if err != nil {
		return nil, err
	}

	return videos, nil
}

// GetVideoBanner implements [EcommerceRepository].
func (d *defaultEcommerceRepository) GetVideoBanner(ctx context.Context, request *dropshipbe.DefaultRequest) (*model.Banner, error) {
	var banner *model.Banner
	cacheKey := constant.VideoBannerKey(request.CountryCode)

	err := d.cache.TakeCtx(ctx, &banner, cacheKey, func(v any) error {
		query := d.db.Model(&model.Banner{}).Where("is_active = ?", true)
		if request.CountryCode != "" {
			query = query.Joins("JOIN countries ON countries.code = banner.country_code").
				Where("countries.code = ?", request.CountryCode)
		}
		return query.First(v).Error
	})

	if err != nil {
		return nil, err
	}

	return banner, nil
}

func (d *defaultEcommerceRepository) GetSliderItems(ctx context.Context, request *dropshipbe.DefaultRequest) ([]model.Slider, error) {
	var sliders []model.Slider
	cacheKey := constant.SliderItemListKey(request.CountryCode)

	err := d.cache.TakeCtx(ctx, &sliders, cacheKey, func(v any) error {
		query := d.db.Model(&model.Slider{}).Where("is_active = ?", true)
		if request.CountryCode != "" {
			query = query.Joins("JOIN countries ON countries.code = sliders.country_code").
				Where("countries.code = ?", request.CountryCode)
		}
		return query.Find(v).Error
	})

	if err != nil {
		return nil, err
	}

	return sliders, nil
}

// GetShop implements [EcommerceRepository].
func (d *defaultEcommerceRepository) GetShop(ctx context.Context, request *dropshipbe.ShopSearchParams) ([]model.Product, error) {
	var products []model.Product
	cacheKey := constant.ShopSearchKey(request.IsFeatured, request.IsNew, request.IsOnSale, request.IsTrending, request.CountryCode)
	err := d.cache.TakeCtx(ctx, &products, cacheKey, func(v any) error {
		query := d.db.Model(&model.Product{}).
			Preload("Country").
			Preload("Categories").
			Preload("Images").
			Preload("PriceTiers").
			Preload("Options.OptionValues").
			Preload("Variants.OptionValues").
			Where("status = ?", "active")

		if request.IsFeatured {
			query = query.Where("is_featured = ?", true)
		}
		if request.IsNew {
			query = query.Where("is_new = ?", true)
		}
		if request.IsOnSale {
			query = query.Where("is_on_sale = ?", true)
		}
		if request.IsTrending {
			query = query.Where("is_trending = ?", true)
		}
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

// GetProductsByCategory implements [EcommerceRepository].
func (d *defaultEcommerceRepository) GetProductsByCategory(ctx context.Context, request *dropshipbe.GetProductsByCategoryRequest) ([]model.Product, error) {
	var products []model.Product
	cacheKey := constant.ProductListByCategoryKey(request.Category, request.CountryCode)
	err := d.cache.TakeCtx(ctx, &products, cacheKey, func(v any) error {
		query := d.db.Model(&model.Product{}).
			Preload("Country").
			Preload("Categories").
			Preload("Images").
			Preload("PriceTiers").
			Preload("Options.OptionValues").
			Preload("Variants.OptionValues").
			Joins("JOIN product_categories ON product_categories.product_id = products.id").
			Joins("JOIN categories ON categories.id = product_categories.category_id").
			Where("categories.slug = ? AND products.status = ?", request.Category, "active")

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

// GetProductReviews implements [EcommerceRepository].
func (d *defaultEcommerceRepository) GetProductReviews(ctx context.Context, request *dropshipbe.GetProductReviewsRequest) ([]model.ProductReview, error) {
	var reviews []model.ProductReview
	cacheKey := constant.ProductReviewListKey(request.Id, request.CountryCode)
	err := d.cache.TakeCtx(ctx, &reviews, cacheKey, func(v any) error {
		query := d.db.Model(&model.ProductReview{}).Where("product_id = ? AND is_verified = ? AND status = ?", request.Id, true, "active")
		return query.Find(v).Error
	})

	if err != nil {
		return nil, err
	}

	return reviews, nil
}

func (d *defaultEcommerceRepository) GetProductFaqs(ctx context.Context, request *dropshipbe.GetProductFaqsRequest) ([]model.ProductFAQ, error) {
	var faqs []model.ProductFAQ
	cacheKey := constant.ProductFaqListKey(request.Id, request.CountryCode)
	err := d.cache.TakeCtx(ctx, &faqs, cacheKey, func(v any) error {
		query := d.db.Model(&model.ProductFAQ{}).Where("product_id = ?", request.Id)
		return query.Find(v).Error
	})

	if err != nil {
		return nil, err
	}

	return faqs, nil
}

// GetProductBySlug implements [EcommerceRepository].
func (d *defaultEcommerceRepository) GetProductBySlug(ctx context.Context, request *dropshipbe.GetProductBySlugRequest) (*model.Product, error) {
	var product model.Product
	cacheKey := constant.ProductDetailKey(request.Slug, request.CountryCode)
	err := d.cache.TakeCtx(ctx, &product, cacheKey, func(v any) error {
		query := d.db.Model(&model.Product{}).
			Preload("Country").
			Preload("Categories").
			Preload("Images").
			Preload("PriceTiers").
			Preload("Options.OptionValues").
			Preload("Variants.OptionValues")

		if request.CountryCode != "" {
			query = query.Joins("JOIN countries ON countries.code = products.country_code").
				Where("countries.code = ?", request.CountryCode)
		}

		return query.Where("slug = ?", request.Slug).First(v).Error
	})

	if err != nil {
		return nil, err
	}

	return &product, nil
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

func (d *defaultEcommerceRepository) GetRelatedProducts(ctx context.Context, request *dropshipbe.GetRelatedProductsRequest) ([]model.Product, error) {
	var products []model.Product

	cacheKey := constant.RelatedProductListKey(request.Id, request.CountryCode)

	err := d.cache.TakeCtx(ctx, &products, cacheKey, func(v any) error {

		query := d.db.Model(&model.Product{}).
			Preload("Country").
			Preload("Categories").
			Preload("Images").
			Preload("PriceTiers").
			Preload("Options.OptionValues").
			Preload("Variants.OptionValues").
			Where("status = ? AND related_product_id = ?", "active", request.Id)

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
