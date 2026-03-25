-- modify "banners" table
ALTER TABLE "banners" ADD COLUMN "country_code" character(2) NULL DEFAULT 'VN', ADD CONSTRAINT "fk_banners_country" FOREIGN KEY ("country_code") REFERENCES "countries" ("code") ON UPDATE NO ACTION ON DELETE NO ACTION;
