-- create "file_metadata" table
CREATE TABLE "file_metadata" (
  "id" bigserial NOT NULL,
  "filename" text NULL,
  "size" integer NULL,
  "content_type" text NULL,
  "file_key" text NULL,
  "uploaded_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- create "countries" table
CREATE TABLE "countries" (
  "code" character(2) NOT NULL,
  "name" character varying(100) NOT NULL,
  "currency" character varying(3) NOT NULL,
  "language_code" character(2) NULL DEFAULT 'vi',
  "is_active" boolean NULL DEFAULT true,
  "created_at" timestamptz NULL,
  PRIMARY KEY ("code")
);
-- create index "idx_country_active" to table: "countries"
CREATE INDEX "idx_country_active" ON "countries" ("is_active") WHERE (is_active = true);
-- create "users" table
CREATE TABLE "users" (
  "id" bigserial NOT NULL,
  "username" text NOT NULL,
  "password" text NOT NULL,
  "role" text NULL,
  "revoke_tokens_before" integer NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_users_username" UNIQUE ("username")
);
-- create "banners" table
CREATE TABLE "banners" (
  "id" bigserial NOT NULL,
  "title" character varying(255) NOT NULL,
  "image_url" character varying(500) NOT NULL,
  "mobile_image_url" character varying(500) NULL,
  "link_url" character varying(500) NULL,
  "position" character varying(50) NOT NULL,
  "sort_order" bigint NULL DEFAULT 0,
  "heading" character varying(255) NULL,
  "sub_heading" text NULL,
  "button_text" character varying(50) NULL,
  "is_active" boolean NULL DEFAULT true,
  "start_date" timestamptz NULL,
  "end_date" timestamptz NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- create index "idx_banner_active" to table: "banners"
CREATE INDEX "idx_banner_active" ON "banners" ("is_active") WHERE (is_active = true);
-- create index "idx_banners_deleted_at" to table: "banners"
CREATE INDEX "idx_banners_deleted_at" ON "banners" ("deleted_at");
-- create "blog_categories" table
CREATE TABLE "blog_categories" (
  "id" bigserial NOT NULL,
  "country_code" character(2) NOT NULL,
  "name" character varying(255) NOT NULL,
  "slug" character varying(255) NOT NULL,
  "description" text NULL,
  "css_class" character varying(100) NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_blog_categories_country" FOREIGN KEY ("country_code") REFERENCES "countries" ("code") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create index "idx_blog_categories_country_code" to table: "blog_categories"
CREATE INDEX "idx_blog_categories_country_code" ON "blog_categories" ("country_code");
-- create index "idx_blog_categories_slug" to table: "blog_categories"
CREATE INDEX "idx_blog_categories_slug" ON "blog_categories" ("slug");
-- create "blog_posts" table
CREATE TABLE "blog_posts" (
  "id" bigserial NOT NULL,
  "country_code" character(2) NOT NULL,
  "category_id" bigint NULL,
  "title" character varying(255) NOT NULL,
  "slug" character varying(255) NOT NULL,
  "excerpt" text NULL,
  "content" text NOT NULL,
  "meta_description" character varying(500) NULL,
  "image_url" character varying(500) NULL,
  "image_alt" character varying(255) NULL,
  "image_width" bigint NULL DEFAULT 0,
  "image_height" bigint NULL DEFAULT 0,
  "author_name" character varying(100) NULL,
  "author_avatar" character varying(500) NULL,
  "tags" jsonb NULL,
  "is_published" boolean NULL DEFAULT false,
  "published_at" timestamptz NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_blog_categories_blog_posts" FOREIGN KEY ("category_id") REFERENCES "blog_categories" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_blog_posts_country" FOREIGN KEY ("country_code") REFERENCES "countries" ("code") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create index "idx_blog_posts_category_id" to table: "blog_posts"
CREATE INDEX "idx_blog_posts_category_id" ON "blog_posts" ("category_id");
-- create index "idx_blog_posts_country_code" to table: "blog_posts"
CREATE INDEX "idx_blog_posts_country_code" ON "blog_posts" ("country_code");
-- create index "idx_blog_posts_deleted_at" to table: "blog_posts"
CREATE INDEX "idx_blog_posts_deleted_at" ON "blog_posts" ("deleted_at");
-- create index "idx_blog_posts_slug" to table: "blog_posts"
CREATE INDEX "idx_blog_posts_slug" ON "blog_posts" ("slug");
-- create index "idx_blog_published" to table: "blog_posts"
CREATE INDEX "idx_blog_published" ON "blog_posts" ("is_published", "published_at");
-- create "categories" table
CREATE TABLE "categories" (
  "id" bigserial NOT NULL,
  "parent_id" bigint NULL,
  "country_code" character(2) NOT NULL,
  "name" character varying(255) NOT NULL,
  "slug" character varying(255) NOT NULL,
  "description" text NULL,
  "image_url" character varying(255) NULL,
  "is_active" boolean NULL DEFAULT true,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_categories_children" FOREIGN KEY ("parent_id") REFERENCES "categories" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_categories_country" FOREIGN KEY ("country_code") REFERENCES "countries" ("code") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create index "idx_cat_country_slug" to table: "categories"
CREATE INDEX "idx_cat_country_slug" ON "categories" ("country_code", "slug");
-- create index "idx_categories_parent_id" to table: "categories"
CREATE INDEX "idx_categories_parent_id" ON "categories" ("parent_id");
-- create index "idx_category_active" to table: "categories"
CREATE INDEX "idx_category_active" ON "categories" ("is_active") WHERE (is_active = true);
-- create "campaigns" table
CREATE TABLE "campaigns" (
  "id" bigserial NOT NULL,
  "name" character varying(255) NOT NULL,
  "description" text NULL,
  "start_date" timestamptz NOT NULL,
  "end_date" timestamptz NOT NULL,
  "is_active" boolean NULL DEFAULT true,
  "created_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- create index "idx_campaigns_deleted_at" to table: "campaigns"
CREATE INDEX "idx_campaigns_deleted_at" ON "campaigns" ("deleted_at");
-- create "coupons" table
CREATE TABLE "coupons" (
  "id" bigserial NOT NULL,
  "campaign_id" bigint NOT NULL,
  "code" character varying(50) NOT NULL,
  "discount_type" character varying(30) NOT NULL,
  "value" numeric(15,2) NOT NULL,
  "min_order_value" numeric(15,2) NULL DEFAULT 0,
  "max_discount_amount" numeric(15,2) NULL,
  "target_type" character varying(30) NULL DEFAULT 'specific_items',
  "usage_limit" bigint NULL DEFAULT 0,
  "usage_limit_per_user" bigint NULL DEFAULT 1,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_campaigns_coupons" FOREIGN KEY ("campaign_id") REFERENCES "campaigns" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create index "idx_coupons_code" to table: "coupons"
CREATE UNIQUE INDEX "idx_coupons_code" ON "coupons" ("code");
-- create index "idx_coupons_deleted_at" to table: "coupons"
CREATE INDEX "idx_coupons_deleted_at" ON "coupons" ("deleted_at");
-- create "products" table
CREATE TABLE "products" (
  "id" bigserial NOT NULL,
  "country_code" character(2) NOT NULL,
  "name" character varying(255) NOT NULL,
  "slug" character varying(255) NOT NULL,
  "metadata" jsonb NULL,
  "description" text NULL,
  "status" character varying(20) NULL DEFAULT 'draft',
  "is_featured" boolean NULL DEFAULT false,
  "is_trending" boolean NULL DEFAULT false,
  "meta_title" character varying(255) NULL,
  "meta_description" character varying(500) NULL,
  "vendor" character varying(100) NULL,
  "product_type" character varying(100) NULL,
  "badge" character varying(50) NULL,
  "sale_label" character varying(50) NULL,
  "sale_tag" character varying(100) NULL,
  "flash_sale_end_time" timestamptz NULL,
  "sold" bigint NULL DEFAULT 0,
  "rating" numeric NULL DEFAULT 0,
  "review_count" bigint NULL DEFAULT 0,
  "tags" jsonb NULL,
  "quantity_enabled" boolean NULL DEFAULT true,
  "quick_shop" boolean NULL DEFAULT true,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- create index "idx_product_country_slug" to table: "products"
CREATE INDEX "idx_product_country_slug" ON "products" ("country_code", "slug");
-- create index "idx_products_status" to table: "products"
CREATE INDEX "idx_products_status" ON "products" ("status");
-- create "variants" table
CREATE TABLE "variants" (
  "id" bigserial NOT NULL,
  "product_id" bigint NOT NULL,
  "sku" character varying(100) NOT NULL,
  "barcode" character varying(100) NULL,
  "price" numeric(15,2) NOT NULL,
  "compare_at_price" numeric(15,2) NULL,
  "cost_price" numeric(15,2) NULL,
  "stock_quantity" bigint NULL DEFAULT 0,
  "weight_gram" bigint NULL DEFAULT 0,
  "image_id" bigint NULL,
  "is_active" boolean NULL DEFAULT true,
  "created_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_products_variants" FOREIGN KEY ("product_id") REFERENCES "products" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "chk_variants_stock_quantity" CHECK (stock_quantity >= 0)
);
-- create index "idx_variants_product_id" to table: "variants"
CREATE INDEX "idx_variants_product_id" ON "variants" ("product_id");
-- create index "idx_variants_sku" to table: "variants"
CREATE UNIQUE INDEX "idx_variants_sku" ON "variants" ("sku");
-- create "coupon_items" table
CREATE TABLE "coupon_items" (
  "coupon_id" bigint NOT NULL,
  "variant_id" bigint NOT NULL,
  PRIMARY KEY ("coupon_id", "variant_id"),
  CONSTRAINT "fk_coupon_items_coupon" FOREIGN KEY ("coupon_id") REFERENCES "coupons" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_coupon_items_variant" FOREIGN KEY ("variant_id") REFERENCES "variants" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create "orders" table
CREATE TABLE "orders" (
  "id" bigserial NOT NULL,
  "order_number" character varying(50) NOT NULL,
  "customer_email" character varying(255) NULL,
  "customer_phone" character varying(20) NULL,
  "total_price" numeric(15,2) NULL DEFAULT 0,
  "subtotal_price" numeric(15,2) NULL DEFAULT 0,
  "total_discounts" numeric(15,2) NULL DEFAULT 0,
  "total_tax" numeric(15,2) NULL DEFAULT 0,
  "shipping_cost" numeric(15,2) NULL DEFAULT 0,
  "currency" character varying(3) NULL DEFAULT 'USD',
  "exchange_rate" numeric(15,6) NULL DEFAULT 1,
  "financial_status" character varying(30) NULL,
  "fulfillment_status" character varying(30) NULL,
  "shipping_address" jsonb NULL,
  "billing_address" jsonb NULL,
  "created_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- create index "idx_orders_deleted_at" to table: "orders"
CREATE INDEX "idx_orders_deleted_at" ON "orders" ("deleted_at");
-- create index "idx_orders_number" to table: "orders"
CREATE UNIQUE INDEX "idx_orders_number" ON "orders" ("order_number");
-- create "coupon_usages" table
CREATE TABLE "coupon_usages" (
  "id" bigserial NOT NULL,
  "order_id" bigint NOT NULL,
  "coupon_id" bigint NOT NULL,
  "customer_email" character varying(255) NOT NULL,
  "discount_amount" numeric(15,2) NOT NULL,
  "used_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_coupon_usages_coupon" FOREIGN KEY ("coupon_id") REFERENCES "coupons" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_orders_coupon_usages" FOREIGN KEY ("order_id") REFERENCES "orders" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create index "idx_coupon_usages_deleted_at" to table: "coupon_usages"
CREATE INDEX "idx_coupon_usages_deleted_at" ON "coupon_usages" ("deleted_at");
-- create "options" table
CREATE TABLE "options" (
  "id" bigserial NOT NULL,
  "product_id" bigint NOT NULL,
  "name" character varying(100) NOT NULL,
  "code" character varying(100) NOT NULL,
  "position" bigint NULL DEFAULT 0,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_products_options" FOREIGN KEY ("product_id") REFERENCES "products" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create index "idx_options_product_id" to table: "options"
CREATE INDEX "idx_options_product_id" ON "options" ("product_id");
-- create "option_values" table
CREATE TABLE "option_values" (
  "id" bigserial NOT NULL,
  "option_id" bigint NOT NULL,
  "value" character varying(100) NOT NULL,
  "color_code" character varying(20) NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_options_option_values" FOREIGN KEY ("option_id") REFERENCES "options" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create index "idx_option_values_option_id" to table: "option_values"
CREATE INDEX "idx_option_values_option_id" ON "option_values" ("option_id");
-- create "policies" table
CREATE TABLE "policies" (
  "id" bigserial NOT NULL,
  "country_code" character(2) NOT NULL,
  "type" character varying(50) NOT NULL,
  "title" character varying(255) NOT NULL,
  "content" text NOT NULL,
  "is_active" boolean NULL DEFAULT true,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_policies_country" FOREIGN KEY ("country_code") REFERENCES "countries" ("code") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create index "idx_policy_country_type" to table: "policies"
CREATE INDEX "idx_policy_country_type" ON "policies" ("country_code", "type");
-- create "product_categories" table
CREATE TABLE "product_categories" (
  "category_id" bigint NOT NULL,
  "product_id" bigint NOT NULL,
  PRIMARY KEY ("category_id", "product_id"),
  CONSTRAINT "fk_product_categories_category" FOREIGN KEY ("category_id") REFERENCES "categories" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_product_categories_product" FOREIGN KEY ("product_id") REFERENCES "products" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create "product_faqs" table
CREATE TABLE "product_faqs" (
  "id" bigserial NOT NULL,
  "product_id" bigint NOT NULL,
  "question" character varying(500) NOT NULL,
  "answer" text NOT NULL,
  "sort_order" bigint NULL DEFAULT 0,
  "created_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_products_fa_qs" FOREIGN KEY ("product_id") REFERENCES "products" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create index "idx_product_faqs_product_id" to table: "product_faqs"
CREATE INDEX "idx_product_faqs_product_id" ON "product_faqs" ("product_id");
-- create "product_images" table
CREATE TABLE "product_images" (
  "id" bigserial NOT NULL,
  "product_id" bigint NOT NULL,
  "image_url" character varying(500) NOT NULL,
  "video_url" character varying(500) NOT NULL,
  "alt_text" character varying(255) NULL,
  "position" bigint NULL DEFAULT 0,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_products_images" FOREIGN KEY ("product_id") REFERENCES "products" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create "suppliers" table
CREATE TABLE "suppliers" (
  "id" bigserial NOT NULL,
  "name" character varying(255) NOT NULL,
  "platform" character varying(50) NOT NULL,
  "homepage_url" character varying(500) NULL,
  "contact_info" character varying(255) NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- create index "idx_suppliers_deleted_at" to table: "suppliers"
CREATE INDEX "idx_suppliers_deleted_at" ON "suppliers" ("deleted_at");
-- create "product_mappings" table
CREATE TABLE "product_mappings" (
  "id" bigserial NOT NULL,
  "local_variant_id" bigint NOT NULL,
  "supplier_id" bigint NOT NULL,
  "source_product_id" character varying(100) NOT NULL,
  "source_variant_id" character varying(100) NULL,
  "source_url" character varying(500) NOT NULL,
  "cost_price_cny" numeric(15,2) NULL DEFAULT 0,
  "cost_price_usd" numeric(15,2) NULL DEFAULT 0,
  "auto_sync_stock" boolean NULL DEFAULT true,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_product_mappings_local_variant" FOREIGN KEY ("local_variant_id") REFERENCES "variants" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_product_mappings_supplier" FOREIGN KEY ("supplier_id") REFERENCES "suppliers" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create index "idx_product_mappings_deleted_at" to table: "product_mappings"
CREATE INDEX "idx_product_mappings_deleted_at" ON "product_mappings" ("deleted_at");
-- create index "uniq_mapping_local_var" to table: "product_mappings"
CREATE UNIQUE INDEX "uniq_mapping_local_var" ON "product_mappings" ("local_variant_id");
-- create "product_price_tiers" table
CREATE TABLE "product_price_tiers" (
  "id" bigserial NOT NULL,
  "product_id" bigint NOT NULL,
  "name" character varying(255) NOT NULL,
  "qty" bigint NOT NULL DEFAULT 1,
  "savings_text" character varying(100) NULL,
  "price" numeric(15,2) NOT NULL,
  "base_price" numeric(15,2) NULL,
  "tag" character varying(50) NULL,
  "tag_class" character varying(100) NULL,
  "created_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_products_price_tiers" FOREIGN KEY ("product_id") REFERENCES "products" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create index "idx_product_price_tiers_deleted_at" to table: "product_price_tiers"
CREATE INDEX "idx_product_price_tiers_deleted_at" ON "product_price_tiers" ("deleted_at");
-- create index "idx_product_price_tiers_product_id" to table: "product_price_tiers"
CREATE INDEX "idx_product_price_tiers_product_id" ON "product_price_tiers" ("product_id");
-- create "product_reviews" table
CREATE TABLE "product_reviews" (
  "id" bigserial NOT NULL,
  "product_id" bigint NOT NULL,
  "author_name" character varying(100) NOT NULL,
  "author_email" character varying(255) NULL,
  "author_avatar" character varying(500) NULL,
  "rating" bigint NOT NULL,
  "content" text NULL,
  "media" jsonb NULL,
  "is_verified" boolean NULL DEFAULT false,
  "status" character varying(20) NULL DEFAULT 'pending',
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_products_reviews" FOREIGN KEY ("product_id") REFERENCES "products" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "chk_product_reviews_rating" CHECK ((rating >= 1) AND (rating <= 5))
);
-- create index "idx_product_reviews_author_email" to table: "product_reviews"
CREATE INDEX "idx_product_reviews_author_email" ON "product_reviews" ("author_email");
-- create index "idx_product_reviews_deleted_at" to table: "product_reviews"
CREATE INDEX "idx_product_reviews_deleted_at" ON "product_reviews" ("deleted_at");
-- create index "idx_product_reviews_product_id" to table: "product_reviews"
CREATE INDEX "idx_product_reviews_product_id" ON "product_reviews" ("product_id");
-- create index "idx_product_reviews_rating" to table: "product_reviews"
CREATE INDEX "idx_product_reviews_rating" ON "product_reviews" ("rating");
-- create index "idx_product_reviews_status" to table: "product_reviews"
CREATE INDEX "idx_product_reviews_status" ON "product_reviews" ("status");
-- create "purchase_orders" table
CREATE TABLE "purchase_orders" (
  "id" bigserial NOT NULL,
  "supplier_id" bigint NOT NULL,
  "platform_order_id" character varying(100) NULL,
  "total_cost" numeric(15,2) NULL,
  "currency" character varying(3) NULL,
  "status" character varying(30) NULL,
  "local_tracking_number" character varying(100) NULL,
  "created_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_purchase_orders_supplier" FOREIGN KEY ("supplier_id") REFERENCES "suppliers" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create index "idx_purchase_orders_deleted_at" to table: "purchase_orders"
CREATE INDEX "idx_purchase_orders_deleted_at" ON "purchase_orders" ("deleted_at");
-- create "purchase_order_items" table
CREATE TABLE "purchase_order_items" (
  "id" bigserial NOT NULL,
  "purchase_order_id" bigint NOT NULL,
  "order_id" bigint NOT NULL,
  "variant_id" bigint NOT NULL,
  "quantity" bigint NOT NULL,
  "cost_per_item" numeric(15,2) NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_purchase_order_items_order" FOREIGN KEY ("order_id") REFERENCES "orders" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_purchase_order_items_variant" FOREIGN KEY ("variant_id") REFERENCES "variants" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_purchase_orders_purchase_order_items" FOREIGN KEY ("purchase_order_id") REFERENCES "purchase_orders" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create index "idx_purchase_order_items_deleted_at" to table: "purchase_order_items"
CREATE INDEX "idx_purchase_order_items_deleted_at" ON "purchase_order_items" ("deleted_at");
-- create "shipments" table
CREATE TABLE "shipments" (
  "id" bigserial NOT NULL,
  "order_id" bigint NOT NULL,
  "purchase_order_id" bigint NULL,
  "tracking_number" character varying(100) NULL,
  "carrier_code" character varying(50) NULL,
  "tracking_url" character varying(500) NULL,
  "status" character varying(30) NULL,
  "shipped_at" timestamptz NULL,
  "estimated_delivery_date" date NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_orders_shipments" FOREIGN KEY ("order_id") REFERENCES "orders" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_purchase_orders_shipments" FOREIGN KEY ("purchase_order_id") REFERENCES "purchase_orders" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create index "idx_shipments_deleted_at" to table: "shipments"
CREATE INDEX "idx_shipments_deleted_at" ON "shipments" ("deleted_at");
-- create "transactions" table
CREATE TABLE "transactions" (
  "id" bigserial NOT NULL,
  "order_id" bigint NOT NULL,
  "gateway" character varying(20) NOT NULL,
  "kind" character varying(20) NULL DEFAULT 'sale',
  "payment_method" character varying(50) NOT NULL,
  "transaction_reference" character varying(255) NULL,
  "amount" numeric(15,2) NOT NULL,
  "currency" character varying(3) NOT NULL,
  "status" character varying(20) NULL DEFAULT 'pending',
  "raw_response" jsonb NULL,
  "error_message" text NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_orders_transactions" FOREIGN KEY ("order_id") REFERENCES "orders" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create index "idx_transactions_deleted_at" to table: "transactions"
CREATE INDEX "idx_transactions_deleted_at" ON "transactions" ("deleted_at");
-- create "variant_value_map" table
CREATE TABLE "variant_value_map" (
  "variant_id" bigint NOT NULL,
  "option_value_id" bigint NOT NULL,
  PRIMARY KEY ("variant_id", "option_value_id"),
  CONSTRAINT "fk_variant_value_map_option_value" FOREIGN KEY ("option_value_id") REFERENCES "option_values" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_variant_value_map_variant" FOREIGN KEY ("variant_id") REFERENCES "variants" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
