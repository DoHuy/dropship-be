-- reverse: create index "idx_products_related_product_id" to table: "products"
DROP INDEX "idx_products_related_product_id";
-- reverse: modify "products" table
ALTER TABLE "products" DROP COLUMN "related_product_id";
