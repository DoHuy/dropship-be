-- reverse: create "variant_value_map" table
DROP TABLE "variant_value_map";
-- reverse: create index "idx_transactions_deleted_at" to table: "transactions"
DROP INDEX "idx_transactions_deleted_at";
-- reverse: create "transactions" table
DROP TABLE "transactions";
-- reverse: create index "idx_shipments_deleted_at" to table: "shipments"
DROP INDEX "idx_shipments_deleted_at";
-- reverse: create "shipments" table
DROP TABLE "shipments";
-- reverse: create index "idx_purchase_order_items_deleted_at" to table: "purchase_order_items"
DROP INDEX "idx_purchase_order_items_deleted_at";
-- reverse: create "purchase_order_items" table
DROP TABLE "purchase_order_items";
-- reverse: create index "idx_purchase_orders_deleted_at" to table: "purchase_orders"
DROP INDEX "idx_purchase_orders_deleted_at";
-- reverse: create "purchase_orders" table
DROP TABLE "purchase_orders";
-- reverse: create index "idx_product_reviews_status" to table: "product_reviews"
DROP INDEX "idx_product_reviews_status";
-- reverse: create index "idx_product_reviews_rating" to table: "product_reviews"
DROP INDEX "idx_product_reviews_rating";
-- reverse: create index "idx_product_reviews_product_id" to table: "product_reviews"
DROP INDEX "idx_product_reviews_product_id";
-- reverse: create index "idx_product_reviews_deleted_at" to table: "product_reviews"
DROP INDEX "idx_product_reviews_deleted_at";
-- reverse: create index "idx_product_reviews_author_email" to table: "product_reviews"
DROP INDEX "idx_product_reviews_author_email";
-- reverse: create "product_reviews" table
DROP TABLE "product_reviews";
-- reverse: create index "idx_product_price_tiers_product_id" to table: "product_price_tiers"
DROP INDEX "idx_product_price_tiers_product_id";
-- reverse: create index "idx_product_price_tiers_deleted_at" to table: "product_price_tiers"
DROP INDEX "idx_product_price_tiers_deleted_at";
-- reverse: create "product_price_tiers" table
DROP TABLE "product_price_tiers";
-- reverse: create index "uniq_mapping_local_var" to table: "product_mappings"
DROP INDEX "uniq_mapping_local_var";
-- reverse: create index "idx_product_mappings_deleted_at" to table: "product_mappings"
DROP INDEX "idx_product_mappings_deleted_at";
-- reverse: create "product_mappings" table
DROP TABLE "product_mappings";
-- reverse: create index "idx_suppliers_deleted_at" to table: "suppliers"
DROP INDEX "idx_suppliers_deleted_at";
-- reverse: create "suppliers" table
DROP TABLE "suppliers";
-- reverse: create "product_images" table
DROP TABLE "product_images";
-- reverse: create index "idx_product_faqs_product_id" to table: "product_faqs"
DROP INDEX "idx_product_faqs_product_id";
-- reverse: create "product_faqs" table
DROP TABLE "product_faqs";
-- reverse: create "product_categories" table
DROP TABLE "product_categories";
-- reverse: create index "idx_policy_country_type" to table: "policies"
DROP INDEX "idx_policy_country_type";
-- reverse: create "policies" table
DROP TABLE "policies";
-- reverse: create index "idx_option_values_option_id" to table: "option_values"
DROP INDEX "idx_option_values_option_id";
-- reverse: create "option_values" table
DROP TABLE "option_values";
-- reverse: create index "idx_options_product_id" to table: "options"
DROP INDEX "idx_options_product_id";
-- reverse: create "options" table
DROP TABLE "options";
-- reverse: create index "idx_coupon_usages_deleted_at" to table: "coupon_usages"
DROP INDEX "idx_coupon_usages_deleted_at";
-- reverse: create "coupon_usages" table
DROP TABLE "coupon_usages";
-- reverse: create index "idx_orders_number" to table: "orders"
DROP INDEX "idx_orders_number";
-- reverse: create index "idx_orders_deleted_at" to table: "orders"
DROP INDEX "idx_orders_deleted_at";
-- reverse: create "orders" table
DROP TABLE "orders";
-- reverse: create "coupon_items" table
DROP TABLE "coupon_items";
-- reverse: create index "idx_variants_sku" to table: "variants"
DROP INDEX "idx_variants_sku";
-- reverse: create index "idx_variants_product_id" to table: "variants"
DROP INDEX "idx_variants_product_id";
-- reverse: create "variants" table
DROP TABLE "variants";
-- reverse: create index "idx_products_status" to table: "products"
DROP INDEX "idx_products_status";
-- reverse: create index "idx_product_country_slug" to table: "products"
DROP INDEX "idx_product_country_slug";
-- reverse: create "products" table
DROP TABLE "products";
-- reverse: create index "idx_coupons_deleted_at" to table: "coupons"
DROP INDEX "idx_coupons_deleted_at";
-- reverse: create index "idx_coupons_code" to table: "coupons"
DROP INDEX "idx_coupons_code";
-- reverse: create "coupons" table
DROP TABLE "coupons";
-- reverse: create index "idx_campaigns_deleted_at" to table: "campaigns"
DROP INDEX "idx_campaigns_deleted_at";
-- reverse: create "campaigns" table
DROP TABLE "campaigns";
-- reverse: create index "idx_category_active" to table: "categories"
DROP INDEX "idx_category_active";
-- reverse: create index "idx_categories_parent_id" to table: "categories"
DROP INDEX "idx_categories_parent_id";
-- reverse: create index "idx_cat_country_slug" to table: "categories"
DROP INDEX "idx_cat_country_slug";
-- reverse: create "categories" table
DROP TABLE "categories";
-- reverse: create index "idx_blog_published" to table: "blog_posts"
DROP INDEX "idx_blog_published";
-- reverse: create index "idx_blog_posts_slug" to table: "blog_posts"
DROP INDEX "idx_blog_posts_slug";
-- reverse: create index "idx_blog_posts_deleted_at" to table: "blog_posts"
DROP INDEX "idx_blog_posts_deleted_at";
-- reverse: create index "idx_blog_posts_country_code" to table: "blog_posts"
DROP INDEX "idx_blog_posts_country_code";
-- reverse: create index "idx_blog_posts_category_id" to table: "blog_posts"
DROP INDEX "idx_blog_posts_category_id";
-- reverse: create "blog_posts" table
DROP TABLE "blog_posts";
-- reverse: create index "idx_blog_categories_slug" to table: "blog_categories"
DROP INDEX "idx_blog_categories_slug";
-- reverse: create index "idx_blog_categories_country_code" to table: "blog_categories"
DROP INDEX "idx_blog_categories_country_code";
-- reverse: create "blog_categories" table
DROP TABLE "blog_categories";
-- reverse: create index "idx_banners_deleted_at" to table: "banners"
DROP INDEX "idx_banners_deleted_at";
-- reverse: create index "idx_banner_active" to table: "banners"
DROP INDEX "idx_banner_active";
-- reverse: create "banners" table
DROP TABLE "banners";
-- reverse: create "users" table
DROP TABLE "users";
-- reverse: create index "idx_country_active" to table: "countries"
DROP INDEX "idx_country_active";
-- reverse: create "countries" table
DROP TABLE "countries";
-- reverse: create "file_metadata" table
DROP TABLE "file_metadata";
