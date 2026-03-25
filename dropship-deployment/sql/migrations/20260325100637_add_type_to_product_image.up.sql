-- modify "product_images" table
ALTER TABLE "product_images" ADD COLUMN "media_type" character varying(50) NULL DEFAULT 'gallery';
