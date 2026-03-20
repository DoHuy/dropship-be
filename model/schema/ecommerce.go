package model

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// ==========================================
// 1. CORE TABLES (Countries, Categories)
// ==========================================

// Country reflects table "countries"
type Country struct {
	Code         string     `gorm:"primaryKey;type:char(2)" json:"code"`
	Name         string     `gorm:"type:varchar(100);not null" json:"name"`
	Currency     string     `gorm:"type:varchar(3);not null" json:"currency"`
	LanguageCode string     `gorm:"type:char(2);default:'vi'" json:"language_code"`
	IsActive     *bool      `gorm:"default:true;index:idx_country_active,where:is_active=true" json:"is_active"` // Pointer để handle DB default true
	CreatedAt    *time.Time `gorm:"type:timestamptz" json:"created_at"`
}

// Category reflects table "categories"
type Category struct {
	ID          uint64  `gorm:"primaryKey;autoIncrement" json:"id"`
	ParentID    *uint64 `gorm:"index" json:"parent_id"`
	CountryCode string  `gorm:"type:char(2);not null;index:idx_cat_country_slug,priority:1" json:"country_code"`

	Name        string `gorm:"type:varchar(255);not null" json:"name"`
	Slug        string `gorm:"type:varchar(255);not null;index:idx_cat_country_slug,priority:2" json:"slug"`
	Description string `gorm:"type:text" json:"description"`
	ImageURL    string `gorm:"type:varchar(255)" json:"image_url"`
	IsActive    *bool  `gorm:"default:true;index:idx_category_active,where:is_active=true" json:"is_active"`

	// Relationships
	Country  *Country   `gorm:"foreignKey:CountryCode;references:Code" json:"country,omitempty"`
	Parent   *Category  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children []Category `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Products []Product  `gorm:"many2many:product_categories;" json:"products,omitempty"`
}

// Product reflects table "products"
type Product struct {
	// --- Định danh cơ bản ---
	ID          uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	CountryCode string `gorm:"type:char(2);not null;index:idx_product_country_slug,priority:1" json:"country_code"`

	Name string `gorm:"type:varchar(255);not null" json:"name"`                                           // Maps to DTO.Name
	Slug string `gorm:"type:varchar(255);not null;index:idx_product_country_slug,priority:2" json:"slug"` // Maps to DTO.Slug

	// --- Metadata & Specs ---
	// Lưu trữ các field: BatteryLife, BluetoothVersion, Weight, WowDelay dưới dạng JSON
	Metadata datatypes.JSON `gorm:"type:jsonb" json:"metadata"` // Maps to DTO.Metadata & DTO.WowDelay

	Description string `gorm:"type:text" json:"description"` // HTML Content

	// Trạng thái: 'draft', 'active', 'archived'
	Status string `gorm:"type:varchar(20);default:'draft';index" json:"status"`

	// --- Marketing & SEO ---
	IsFeatured      bool   `gorm:"default:false" json:"is_featured"`
	IsTrending      bool   `gorm:"default:false" json:"is_trending"` // Maps to DTO.IsTrending
	MetaTitle       string `gorm:"type:varchar(255)" json:"meta_title"`
	MetaDescription string `gorm:"type:varchar(500)" json:"meta_description"`
	Vendor          string `gorm:"type:varchar(100)" json:"vendor"`
	ProductType     string `gorm:"type:varchar(100)" json:"product_type"`

	// --- Labels & Badges ---
	Badge     *string `gorm:"type:varchar(50)" json:"badge"`      // DTO.Badge ("Best Seller")
	SaleLabel *string `gorm:"type:varchar(50)" json:"sale_label"` // DTO.SaleLabel ("Hot")
	SaleTag   *string `gorm:"type:varchar(100)" json:"sale_tag"`  // DTO.SaleTag ("Save $50")

	// --- Flash Sale ---
	// Dùng để tính toán DTO.FlashSaleEnd & DTO.IsOutOfSale
	FlashSaleEndTime *time.Time `gorm:"type:timestamptz" json:"flash_sale_end_time"`

	// --- Thống kê (Cache fields) ---
	Sold        int     `gorm:"default:0" json:"sold"`
	Rating      float64 `gorm:"default:0" json:"rating"`
	ReviewCount int     `gorm:"default:0" json:"review_count"`

	// --- Cấu hình hiển thị ---
	Tags            datatypes.JSON `gorm:"type:jsonb" json:"tags"` // Maps to DTO.Tags ([]string)
	QuantityEnabled bool           `gorm:"default:true" json:"quantity_enabled"`
	QuickShop       bool           `gorm:"default:true" json:"quick_shop"`

	CreatedAt *time.Time `gorm:"type:timestamptz" json:"created_at"`
	UpdatedAt *time.Time `gorm:"type:timestamptz" json:"updated_at"`

	// --- Relationships ---
	Categories []Category `gorm:"many2many:product_categories;" json:"categories,omitempty"`

	// Images gộp chung cả Gallery, DescriptionImages và SocialVideos (phân loại bằng Type)
	Images []ProductImage `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE" json:"images,omitempty"`

	Options    []Option           `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE" json:"options,omitempty"`
	Variants   []Variant          `gorm:"foreignKey:ProductID" json:"variants,omitempty"`
	PriceTiers []ProductPriceTier `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE" json:"price_tiers,omitempty"`
	Reviews    []ProductReview    `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE" json:"reviews,omitempty"`
	FAQs       []ProductFAQ       `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE" json:"faqs,omitempty"` // [NEW]
}

// ProductFAQ reflects table "product_faqs" (NEW TABLE)
// Maps to ProductFAQDTO
type ProductFAQ struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductID uint64 `gorm:"not null;index" json:"product_id"`

	Question  string `gorm:"type:varchar(500);not null" json:"question"`
	Answer    string `gorm:"type:text;not null" json:"answer"`
	SortOrder int    `gorm:"default:0" json:"sort_order"`

	CreatedAt *time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// ProductImage reflects table "product_images"
type ProductImage struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductID uint64 `gorm:"not null" json:"product_id"`
	ImageURL  string `gorm:"type:varchar(500);not null" json:"image_url"`
	VideoURL  string `gorm:"type:varchar(500);not null" json:"video_url"`
	AltText   string `gorm:"type:varchar(255)" json:"alt_text"`
	Position  int    `gorm:"default:0" json:"position"`
}

// ProductPriceTier reflects table "product_price_tiers" (NEW MODEL)
// Handles "Buy 1", "Buy 2 Save X" logic
type ProductPriceTier struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductID uint64 `gorm:"not null;index" json:"product_id"`

	Name string `gorm:"type:varchar(255);not null" json:"name"` // e.g. "Buy 1 item"
	Qty  int    `gorm:"not null;default:1" json:"qty"`          // e.g. 1

	SavingsText string `gorm:"type:varchar(100)" json:"savings"` // e.g. "Standard Price"

	// Price stored as numeric(15,2). Frontend should format currency ($125)
	Price     float64  `gorm:"type:numeric(15,2);not null" json:"new_price"` // e.g. 125.00
	BasePrice *float64 `gorm:"type:numeric(15,2)" json:"old_price"`          // e.g. 130.00

	Tag      *string `gorm:"type:varchar(50)" json:"tag"`        // e.g. "Best Value"
	TagClass string  `gorm:"type:varchar(100)" json:"tag_class"` // e.g. "badge-warning"

	CreatedAt *time.Time     `gorm:"autoCreateTime;type:timestamptz" json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index;type:timestamptz" json:"deleted_at"`

	// Relationships
	Product *Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

// Option reflects table "options"
type Option struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductID uint64 `gorm:"not null;index" json:"product_id"`

	Name     string `gorm:"type:varchar(100);not null" json:"name"` // "Color"
	Code     string `gorm:"type:varchar(100);not null" json:"code"` // "color" (backend slug)
	Position int    `gorm:"default:0" json:"position"`

	OptionValues []OptionValue `gorm:"foreignKey:OptionID;constraint:OnDelete:CASCADE" json:"values,omitempty"`
}

// OptionValue reflects table "option_values"
type OptionValue struct {
	ID       uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	OptionID uint64 `gorm:"not null;index" json:"option_id"`

	Value     string `gorm:"type:varchar(100);not null" json:"value"` // "Matte Black"
	ColorCode string `gorm:"type:varchar(20)" json:"color_code"`      // "#000000"
}

// Variant reflects table "variants"
type Variant struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductID uint64 `gorm:"not null;index" json:"product_id"`

	Sku     string `gorm:"type:varchar(100);not null;uniqueIndex:idx_variants_sku" json:"sku"`
	Barcode string `gorm:"type:varchar(100)" json:"barcode"`

	// Pricing
	Price          float64 `gorm:"type:numeric(15,2);not null" json:"price"`
	CompareAtPrice float64 `gorm:"type:numeric(15,2)" json:"compare_at_price"`
	CostPrice      float64 `gorm:"type:numeric(15,2)" json:"cost_price"`

	// Inventory & Specs
	StockQuantity int    `gorm:"default:0;check:stock_quantity >= 0" json:"stock_quantity"`
	WeightGram    int    `gorm:"default:0" json:"weight_gram"`
	ImageID       *int64 `gorm:"default:null" json:"image_id"` // Link tới ProductImage

	IsActive  *bool      `gorm:"default:true" json:"is_active"`
	CreatedAt *time.Time `gorm:"type:timestamptz" json:"created_at"`

	// Relationships
	// Many-to-Many để map Variant với OptionValues (VD: Variant A = [Red, XL])
	OptionValues []OptionValue `gorm:"many2many:variant_value_map;" json:"option_values,omitempty"`
}

// ==========================================
// 3. MARKETING
// ==========================================

// Campaign reflects table "campaigns"
type Campaign struct {
	ID          uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string         `gorm:"type:varchar(255);not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	StartDate   time.Time      `gorm:"not null;type:timestamptz" json:"start_date"`
	EndDate     time.Time      `gorm:"not null;type:timestamptz" json:"end_date"`
	IsActive    *bool          `gorm:"default:true" json:"is_active"`
	CreatedAt   *time.Time     `gorm:"autoCreateTime;type:timestamptz" json:"created_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index;type:timestamptz" json:"deleted_at"`

	Coupons []Coupon `gorm:"foreignKey:CampaignID;constraint:OnDelete:CASCADE" json:"coupons,omitempty"`
}

// Coupon reflects table "coupons"
type Coupon struct {
	ID         uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	CampaignID uint64 `gorm:"not null" json:"campaign_id"`
	Code       string `gorm:"type:varchar(50);not null;uniqueIndex:idx_coupons_code" json:"code"`

	DiscountType string  `gorm:"type:varchar(30);not null" json:"discount_type"`
	Value        float64 `gorm:"type:numeric(15,2);not null" json:"value"`

	MinOrderValue     float64 `gorm:"type:numeric(15,2);default:0" json:"min_order_value"`
	MaxDiscountAmount float64 `gorm:"type:numeric(15,2)" json:"max_discount_amount"`

	TargetType string `gorm:"type:varchar(30);default:'specific_items'" json:"target_type"`

	UsageLimit        int `gorm:"default:0" json:"usage_limit"`
	UsageLimitPerUser int `gorm:"default:1" json:"usage_limit_per_user"`

	DeletedAt gorm.DeletedAt `gorm:"index;type:timestamptz" json:"deleted_at"`

	Campaign           *Campaign `gorm:"foreignKey:CampaignID" json:"campaign,omitempty"`
	ApplicableVariants []Variant `gorm:"many2many:coupon_items;" json:"applicable_variants,omitempty"`
}

// Banner reflects table "banners"
type Banner struct {
	ID    uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	Title string `gorm:"type:varchar(255);not null" json:"title"`

	ImageURL       string `gorm:"type:varchar(500);not null" json:"image_url"`
	MobileImageURL string `gorm:"type:varchar(500)" json:"mobile_image_url"`
	LinkURL        string `gorm:"type:varchar(500)" json:"link_url"`

	Position   string `gorm:"type:varchar(50);not null" json:"position"`
	SortOrder  int    `gorm:"default:0" json:"sort_order"`
	Heading    string `gorm:"type:varchar(255)" json:"heading"`
	SubHeading string `gorm:"type:text" json:"sub_heading"`
	ButtonText string `gorm:"type:varchar(50)" json:"button_text"`

	IsActive *bool `gorm:"default:true;index:idx_banner_active,where:is_active=true" json:"is_active"`

	StartDate *time.Time `gorm:"type:timestamptz" json:"start_date"`
	EndDate   *time.Time `gorm:"type:timestamptz" json:"end_date"`

	CreatedAt *time.Time     `gorm:"autoCreateTime;type:timestamptz" json:"created_at"`
	UpdatedAt *time.Time     `gorm:"autoUpdateTime;type:timestamptz" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index;type:timestamptz" json:"deleted_at"`
}

// ==========================================
// 4. ORDER & FULFILLMENT
// ==========================================

// Order reflects table "orders"
type Order struct {
	ID            uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderNumber   string `gorm:"type:varchar(50);not null;uniqueIndex:idx_orders_number" json:"order_number"`
	CustomerEmail string `gorm:"type:varchar(255)" json:"customer_email"`
	CustomerPhone string `gorm:"type:varchar(20)" json:"customer_phone"`

	TotalPrice     float64 `gorm:"type:numeric(15,2);default:0" json:"total_price"`
	SubtotalPrice  float64 `gorm:"type:numeric(15,2);default:0" json:"subtotal_price"`
	TotalDiscounts float64 `gorm:"type:numeric(15,2);default:0" json:"total_discounts"`
	TotalTax       float64 `gorm:"type:numeric(15,2);default:0" json:"total_tax"`
	ShippingCost   float64 `gorm:"type:numeric(15,2);default:0" json:"shipping_cost"`

	Currency     string  `gorm:"type:varchar(3);default:'USD'" json:"currency"`
	ExchangeRate float64 `gorm:"type:numeric(15,6);default:1.0" json:"exchange_rate"`

	FinancialStatus   string `gorm:"type:varchar(30)" json:"financial_status"`
	FulfillmentStatus string `gorm:"type:varchar(30)" json:"fulfillment_status"`

	// JSONB columns with GIN index (via SQL)
	ShippingAddress datatypes.JSON `gorm:"type:jsonb" json:"shipping_address"`
	BillingAddress  datatypes.JSON `gorm:"type:jsonb" json:"billing_address"`

	CreatedAt *time.Time     `gorm:"autoCreateTime;type:timestamptz" json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index;type:timestamptz" json:"deleted_at"`

	Shipments    []Shipment    `gorm:"foreignKey:OrderID" json:"shipments,omitempty"`
	CouponUsages []CouponUsage `gorm:"foreignKey:OrderID" json:"coupon_usages,omitempty"`
	Transactions []Transaction `gorm:"foreignKey:OrderID" json:"transactions,omitempty"`
}

// CouponUsage reflects table "coupon_usages"
type CouponUsage struct {
	ID             uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderID        uint64         `gorm:"not null" json:"order_id"`
	CouponID       uint64         `gorm:"not null" json:"coupon_id"`
	CustomerEmail  string         `gorm:"type:varchar(255);not null" json:"customer_email"`
	DiscountAmount float64        `gorm:"type:numeric(15,2);not null" json:"discount_amount"`
	UsedAt         *time.Time     `gorm:"autoCreateTime;type:timestamptz" json:"used_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index;type:timestamptz" json:"deleted_at"`

	Order  *Order  `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	Coupon *Coupon `gorm:"foreignKey:CouponID" json:"coupon,omitempty"`
}

// Transaction reflects table "transactions"
type Transaction struct {
	ID      uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderID uint64 `gorm:"not null" json:"order_id"`

	Gateway       string `gorm:"type:varchar(20);not null" json:"gateway"`
	Kind          string `gorm:"type:varchar(20);default:'sale'" json:"kind"`
	PaymentMethod string `gorm:"type:varchar(50);not null" json:"payment_method"`

	TransactionReference string `gorm:"type:varchar(255)" json:"transaction_reference"`

	Amount   float64 `gorm:"type:numeric(15,2);not null" json:"amount"`
	Currency string  `gorm:"type:varchar(3);not null" json:"currency"`
	Status   string  `gorm:"type:varchar(20);default:'pending'" json:"status"`

	RawResponse  datatypes.JSON `gorm:"type:jsonb" json:"raw_response"`
	ErrorMessage string         `gorm:"type:text" json:"error_message"`

	CreatedAt *time.Time     `gorm:"autoCreateTime;type:timestamptz" json:"created_at"`
	UpdatedAt *time.Time     `gorm:"autoUpdateTime;type:timestamptz" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index;type:timestamptz" json:"deleted_at"`
}

// ==========================================
// 5. SUPPLY CHAIN
// ==========================================

// Supplier reflects table "suppliers"
type Supplier struct {
	ID          uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string         `gorm:"type:varchar(255);not null" json:"name"`
	Platform    string         `gorm:"type:varchar(50);not null" json:"platform"`
	HomepageURL string         `gorm:"type:varchar(500)" json:"homepage_url"`
	ContactInfo string         `gorm:"type:varchar(255)" json:"contact_info"`
	DeletedAt   gorm.DeletedAt `gorm:"index;type:timestamptz" json:"deleted_at"`
}

// ProductMapping reflects table "product_mappings"
type ProductMapping struct {
	ID             uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	LocalVariantID uint64 `gorm:"not null;uniqueIndex:uniq_mapping_local_var" json:"local_variant_id"`
	SupplierID     uint64 `gorm:"not null" json:"supplier_id"`

	SourceProductID string `gorm:"type:varchar(100);not null" json:"source_product_id"`
	SourceVariantID string `gorm:"type:varchar(100)" json:"source_variant_id"`
	SourceURL       string `gorm:"type:varchar(500);not null" json:"source_url"`

	CostPriceCNY  float64 `gorm:"type:numeric(15,2);default:0" json:"cost_price_cny"`
	CostPriceUSD  float64 `gorm:"type:numeric(15,2);default:0" json:"cost_price_usd"`
	AutoSyncStock *bool   `gorm:"default:true" json:"auto_sync_stock"`

	DeletedAt gorm.DeletedAt `gorm:"index;type:timestamptz" json:"deleted_at"`

	LocalVariant *Variant  `gorm:"foreignKey:LocalVariantID" json:"local_variant,omitempty"`
	Supplier     *Supplier `gorm:"foreignKey:SupplierID" json:"supplier,omitempty"`
}

// PurchaseOrder reflects table "purchase_orders"
type PurchaseOrder struct {
	ID              uint64  `gorm:"primaryKey;autoIncrement" json:"id"`
	SupplierID      uint64  `gorm:"not null" json:"supplier_id"`
	PlatformOrderID string  `gorm:"type:varchar(100)" json:"platform_order_id"`
	TotalCost       float64 `gorm:"type:numeric(15,2)" json:"total_cost"`
	Currency        string  `gorm:"type:varchar(3)" json:"currency"`
	Status          string  `gorm:"type:varchar(30)" json:"status"`

	LocalTrackingNumber string         `gorm:"type:varchar(100)" json:"local_tracking_number"`
	CreatedAt           *time.Time     `gorm:"autoCreateTime;type:timestamptz" json:"created_at"`
	DeletedAt           gorm.DeletedAt `gorm:"index;type:timestamptz" json:"deleted_at"`

	Supplier           *Supplier           `gorm:"foreignKey:SupplierID" json:"supplier,omitempty"`
	PurchaseOrderItems []PurchaseOrderItem `gorm:"foreignKey:PurchaseOrderID" json:"items,omitempty"`
	Shipments          []Shipment          `gorm:"foreignKey:PurchaseOrderID" json:"shipments,omitempty"`
}

// Shipment reflects table "shipments"
type Shipment struct {
	ID              uint64  `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderID         uint64  `gorm:"not null" json:"order_id"`
	PurchaseOrderID *uint64 `gorm:"default:null" json:"purchase_order_id"`

	TrackingNumber string `gorm:"type:varchar(100)" json:"tracking_number"`
	CarrierCode    string `gorm:"type:varchar(50)" json:"carrier_code"`
	TrackingURL    string `gorm:"type:varchar(500)" json:"tracking_url"`
	Status         string `gorm:"type:varchar(30)" json:"status"`

	ShippedAt             *time.Time     `gorm:"type:timestamptz" json:"shipped_at"`
	EstimatedDeliveryDate *time.Time     `gorm:"type:date" json:"estimated_delivery_date"` // type:date trong postgres
	DeletedAt             gorm.DeletedAt `gorm:"index;type:timestamptz" json:"deleted_at"`

	Order         *Order         `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	PurchaseOrder *PurchaseOrder `gorm:"foreignKey:PurchaseOrderID" json:"purchase_order,omitempty"`
}

// PurchaseOrderItem reflects table "purchase_order_items"
type PurchaseOrderItem struct {
	ID              uint64  `gorm:"primaryKey;autoIncrement" json:"id"`
	PurchaseOrderID uint64  `gorm:"not null" json:"purchase_order_id"`
	OrderID         uint64  `gorm:"not null" json:"order_id"`
	VariantID       uint64  `gorm:"not null" json:"variant_id"`
	Quantity        int     `gorm:"not null" json:"quantity"`
	CostPerItem     float64 `gorm:"type:numeric(15,2)" json:"cost_per_item"`

	DeletedAt gorm.DeletedAt `gorm:"index;type:timestamptz" json:"deleted_at"`

	// Relationships don't always need to be loaded, but good to define
	Order   *Order   `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	Variant *Variant `gorm:"foreignKey:VariantID" json:"variant,omitempty"`
}

// ==========================================
// 6. CONTENT & SUPPORT
// ==========================================

// Policy reflects table "policies" (Privacy, Terms, Refunds, etc.)
type Policy struct {
	ID          uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	CountryCode string `gorm:"type:char(2);not null;index:idx_policy_country_type,priority:1" json:"country_code"`

	// Type suggestions: 'privacy_policy', 'terms_conditions', 'return_refund'
	Type  string `gorm:"type:varchar(50);not null;index:idx_policy_country_type,priority:2" json:"type"`
	Title string `gorm:"type:varchar(255);not null" json:"title"`

	Content string `gorm:"type:text;not null" json:"content"` // HTML or Markdown content

	IsActive  *bool      `gorm:"default:true" json:"is_active"`
	CreatedAt *time.Time `gorm:"autoCreateTime;type:timestamptz" json:"created_at"`
	UpdatedAt *time.Time `gorm:"autoUpdateTime;type:timestamptz" json:"updated_at"`

	// Relationships
	Country *Country `gorm:"foreignKey:CountryCode;references:Code" json:"country,omitempty"`
}

// BlogCategory: Phản ánh bảng "blog_categories"
// Đã thêm trường CssClass từ yêu cầu input 2 để hỗ trợ styling.
type BlogCategory struct {
	ID          uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	CountryCode string `gorm:"type:char(2);not null;index" json:"country_code"`

	Name        string `gorm:"type:varchar(255);not null" json:"name"`
	Slug        string `gorm:"type:varchar(255);not null;index" json:"slug"`
	Description string `gorm:"type:text" json:"description"`

	// Field mới từ Input 2: Class CSS để style màu sắc cho danh mục (e.g., "bg-blue-500")
	CssClass string `gorm:"type:varchar(100)" json:"css_class"`

	// Quan hệ
	Country   *Country   `gorm:"foreignKey:CountryCode;references:Code" json:"country,omitempty"`
	BlogPosts []BlogPost `gorm:"foreignKey:CategoryID" json:"blog_posts,omitempty"`
}

// BlogPost: Phản ánh bảng "blog_posts"
// Hợp nhất các biến thể hình ảnh và tác giả vào các cột chuẩn.
type BlogPost struct {
	// --- Định danh & Quan hệ ---
	ID          uint64  `gorm:"primaryKey;autoIncrement" json:"id"`
	CountryCode string  `gorm:"type:char(2);not null;index" json:"country_code"`
	CategoryID  *uint64 `gorm:"index" json:"category_id"` // Optional category

	// --- Nội dung chính ---
	Title string `gorm:"type:varchar(255);not null" json:"title"`
	Slug  string `gorm:"type:varchar(255);not null;index" json:"slug"`

	// Excerpt: Dùng chung cho cả "Excerpt", "Desc" (ngắn gọn)
	Excerpt string `gorm:"type:text" json:"excerpt"`

	// Content: Nội dung HTML đầy đủ
	Content string `gorm:"type:text;not null" json:"content"`

	// MetaDescription: Dùng cho trường "Description" đầy đủ (SEO)
	MetaDescription string `gorm:"type:varchar(500)" json:"meta_description"`

	// --- Hình ảnh (Đã chuẩn hóa) ---
	// Thay vì lưu ImgSrc, ImageSrc, Img... ta chỉ dùng ImageURL làm chuẩn
	ImageURL    string `gorm:"type:varchar(500)" json:"image_url"`
	ImageAlt    string `gorm:"type:varchar(255)" json:"image_alt"` // SEO/Accessibility
	ImageWidth  int    `gorm:"default:0" json:"image_width"`       // Chiều rộng ảnh
	ImageHeight int    `gorm:"default:0" json:"image_height"`      // Chiều cao ảnh

	// --- Thông tin tác giả ---
	// AuthorName thay thế cho field "Author" chung chung để rõ nghĩa hơn
	AuthorName   string `gorm:"type:varchar(100)" json:"author_name"`
	AuthorAvatar string `gorm:"type:varchar(500)" json:"author_avatar"` // Ảnh đại diện tác giả

	// --- Metadata & Phân loại ---
	// Sử dụng JSONB của Postgres để lưu mảng Tags, linh hoạt hơn []string thuần
	Tags datatypes.JSON `gorm:"type:jsonb" json:"tags"` // e.g. ["news", "sale"]

	// --- Trạng thái & Thời gian ---
	IsPublished bool       `gorm:"default:false;index:idx_blog_published" json:"is_published"`
	PublishedAt *time.Time `gorm:"index:idx_blog_published" json:"published_at"` // Tương ứng field "Date"

	CreatedAt *time.Time     `gorm:"autoCreateTime;type:timestamptz" json:"created_at"`
	UpdatedAt *time.Time     `gorm:"autoUpdateTime;type:timestamptz" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index;type:timestamptz" json:"deleted_at"`

	// --- Quan hệ (Associations) ---
	Country  *Country      `gorm:"foreignKey:CountryCode;references:Code" json:"country,omitempty"`
	Category *BlogCategory `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}

// ==========================================
// 7. INVENTORY MANAGEMENT
// ==========================================

// Inventory reflects table "inventories"
// Tracks current stock levels for variants across different warehouses/locations
type Inventory struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	VariantID uint64 `gorm:"not null;uniqueIndex:idx_inv_variant_location,priority:1" json:"variant_id"`
	Location  string `gorm:"type:varchar(100);not null;uniqueIndex:idx_inv_variant_location,priority:2;default:'main'" json:"location"` // e.g. 'main', 'warehouse_a', 'backup'

	QuantityOnHand    int `gorm:"not null;default:0;check:quantity_on_hand >= 0" json:"quantity_on_hand"`                           // Current stock
	QuantityReserved  int `gorm:"not null;default:0;check:quantity_reserved >= 0" json:"quantity_reserved"`                         // Reserved for orders
	QuantityAvailable int `gorm:"type:GENERATED ALWAYS AS (quantity_on_hand - quantity_reserved) STORED" json:"quantity_available"` // Available = OnHand - Reserved

	MinStockLevel   int        `gorm:"default:10" json:"min_stock_level"`  // Reorder point
	ReorderQuantity int        `gorm:"default:50" json:"reorder_quantity"` // Qty to order when below min
	LastCountedAt   *time.Time `gorm:"type:timestamptz" json:"last_counted_at"`

	CreatedAt *time.Time     `gorm:"autoCreateTime;type:timestamptz" json:"created_at"`
	UpdatedAt *time.Time     `gorm:"autoUpdateTime;type:timestamptz" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index;type:timestamptz" json:"deleted_at"`

	// Relationships
	Variant               *Variant               `gorm:"foreignKey:VariantID" json:"variant,omitempty"`
	InventoryTransactions []InventoryTransaction `gorm:"foreignKey:InventoryID" json:"transactions,omitempty"`
}

// InventoryTransaction reflects table "inventory_transactions"
// Audit trail for all inventory movements
type InventoryTransaction struct {
	ID          uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	InventoryID uint64 `gorm:"not null;index" json:"inventory_id"`
	VariantID   uint64 `gorm:"not null;index" json:"variant_id"`

	TransactionType string `gorm:"type:varchar(50);not null;index" json:"transaction_type"` // 'stock_in', 'stock_out', 'adjustment', 'damage', 'return', 'reservation', 'release'

	QuantityChange int    `gorm:"not null" json:"quantity_change"`        // Can be positive or negative
	Reference      string `gorm:"type:varchar(100)" json:"reference"`     // e.g. PurchaseOrderID, OrderID, AdjustmentID
	ReferenceType  string `gorm:"type:varchar(50)" json:"reference_type"` // e.g. 'purchase_order', 'order', 'manual_adjustment'
	Notes          string `gorm:"type:text" json:"notes"`                 // Reason for transaction

	CreatedBy string         `gorm:"type:varchar(255)" json:"created_by"` // User who made the change
	CreatedAt *time.Time     `gorm:"autoCreateTime;type:timestamptz" json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index;type:timestamptz" json:"deleted_at"`

	// Relationships
	Inventory *Inventory `gorm:"foreignKey:InventoryID" json:"inventory,omitempty"`
	Variant   *Variant   `gorm:"foreignKey:VariantID" json:"variant,omitempty"`
}

// ProductReview reflects table "product_reviews"
// Maps to ReviewItemDTO & ReviewAggregateDTO
type ProductReview struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductID uint64 `gorm:"not null;index" json:"product_id"`

	// --- User Info ---
	AuthorName   string `gorm:"type:varchar(100);not null" json:"name"`
	AuthorEmail  string `gorm:"type:varchar(255);index" json:"email"`
	AuthorAvatar string `gorm:"type:varchar(500)" json:"avatar"`

	// --- Rating & Content ---
	Rating  int    `gorm:"not null;check:rating >= 1 AND rating <= 5;index" json:"rating"`
	Content string `gorm:"type:text" json:"content"` // Maps to DTO.Comment

	// --- Media (Multi Images & Videos) ---
	// Lưu cấu trúc JSON: {"images": ["url1", "url2"], "videos": ["url3"]}
	Media datatypes.JSON `gorm:"type:jsonb" json:"media"`

	// --- Status ---
	IsVerified bool   `gorm:"default:false" json:"is_verified"` // Maps to DTO.Verified
	Status     string `gorm:"type:varchar(20);default:'pending';index" json:"status"`

	CreatedAt time.Time      `gorm:"autoCreateTime;type:timestamptz" json:"created_at"` // Maps to DTO.Date
	UpdatedAt *time.Time     `gorm:"autoUpdateTime;type:timestamptz" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index;type:timestamptz" json:"deleted_at"`
}

type FileMetadata struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Filename    string    `gorm:"type:text" json:"filename"`
	Size        int64     `gorm:"type:integer" json:"size"`
	ContentType string    `gorm:"type:text" json:"content_type"`
	FileKey     string    `gorm:"type:text" json:"file_key"`
	UploadedAt  time.Time `gorm:"autoCreateTime" json:"uploaded_at"`
}

func (FileMetadata) TableName() string { return "file_metadata" }
