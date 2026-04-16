-- create "frequently_boughts" table
CREATE TABLE "frequently_boughts" (
  "id" bigserial NOT NULL,
  "product_id" bigint NOT NULL,
  "bought_with_product_id" bigint NOT NULL,
  "sort_order" bigint NULL DEFAULT 0,
  "is_active" boolean NULL DEFAULT true,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_frequently_boughts_bought_with_product" FOREIGN KEY ("bought_with_product_id") REFERENCES "products" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_products_frequently_bought" FOREIGN KEY ("product_id") REFERENCES "products" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create index "idx_product_bought_with" to table: "frequently_boughts"
CREATE UNIQUE INDEX "idx_product_bought_with" ON "frequently_boughts" ("product_id", "bought_with_product_id");
