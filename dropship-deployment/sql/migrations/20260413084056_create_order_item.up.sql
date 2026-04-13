-- create "order_items" table
CREATE TABLE "order_items" (
  "id" bigserial NOT NULL,
  "order_id" bigint NOT NULL,
  "variant_id" bigint NOT NULL,
  "product_id" bigint NOT NULL,
  "product_name" character varying(255) NOT NULL,
  "variant_name" character varying(255) NULL,
  "sku" character varying(100) NULL,
  "quantity" bigint NOT NULL,
  "price" numeric(15,2) NOT NULL,
  "total" numeric(15,2) NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_order_items_order" FOREIGN KEY ("order_id") REFERENCES "orders" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_order_items_product" FOREIGN KEY ("product_id") REFERENCES "products" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_order_items_variant" FOREIGN KEY ("variant_id") REFERENCES "variants" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "chk_order_items_quantity" CHECK (quantity > 0)
);
-- create index "idx_order_items_deleted_at" to table: "order_items"
CREATE INDEX "idx_order_items_deleted_at" ON "order_items" ("deleted_at");
-- create index "idx_order_items_order_id" to table: "order_items"
CREATE INDEX "idx_order_items_order_id" ON "order_items" ("order_id");
-- create index "idx_order_items_product_id" to table: "order_items"
CREATE INDEX "idx_order_items_product_id" ON "order_items" ("product_id");
-- create index "idx_order_items_variant_id" to table: "order_items"
CREATE INDEX "idx_order_items_variant_id" ON "order_items" ("variant_id");
