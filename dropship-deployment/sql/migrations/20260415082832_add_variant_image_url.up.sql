-- modify "variants" table
ALTER TABLE "variants" ADD COLUMN "image_url" character varying(255) NULL;
-- modify "order_items" table
ALTER TABLE "order_items" DROP CONSTRAINT "fk_order_items_order", ADD CONSTRAINT "fk_orders_order_items" FOREIGN KEY ("order_id") REFERENCES "orders" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
