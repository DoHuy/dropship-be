-- modify "products" table
ALTER TABLE "products" ADD COLUMN "related_product_id" bigint NULL;
-- create index "idx_products_related_product_id" to table: "products"
CREATE INDEX "idx_products_related_product_id" ON "products" ("related_product_id");
