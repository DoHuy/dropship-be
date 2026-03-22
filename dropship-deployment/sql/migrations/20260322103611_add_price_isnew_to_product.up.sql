-- modify "products" table
ALTER TABLE "products" ADD COLUMN "price" numeric(15,2) NOT NULL, ADD COLUMN "is_new" boolean NULL DEFAULT false;
