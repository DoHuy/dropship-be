-- reverse: modify "order_items" table
ALTER TABLE "order_items" DROP CONSTRAINT "fk_orders_order_items", ADD CONSTRAINT "fk_order_items_order" FOREIGN KEY ("order_id") REFERENCES "orders" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- reverse: modify "variants" table
ALTER TABLE "variants" DROP COLUMN "image_url";
