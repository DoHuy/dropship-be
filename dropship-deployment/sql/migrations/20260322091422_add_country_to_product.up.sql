-- modify "products" table
ALTER TABLE "products" ADD CONSTRAINT "fk_products_country" FOREIGN KEY ("country_code") REFERENCES "countries" ("code") ON UPDATE NO ACTION ON DELETE NO ACTION;
