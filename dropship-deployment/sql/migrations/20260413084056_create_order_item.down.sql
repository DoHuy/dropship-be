-- reverse: create index "idx_order_items_variant_id" to table: "order_items"
DROP INDEX "idx_order_items_variant_id";
-- reverse: create index "idx_order_items_product_id" to table: "order_items"
DROP INDEX "idx_order_items_product_id";
-- reverse: create index "idx_order_items_order_id" to table: "order_items"
DROP INDEX "idx_order_items_order_id";
-- reverse: create index "idx_order_items_deleted_at" to table: "order_items"
DROP INDEX "idx_order_items_deleted_at";
-- reverse: create "order_items" table
DROP TABLE "order_items";
