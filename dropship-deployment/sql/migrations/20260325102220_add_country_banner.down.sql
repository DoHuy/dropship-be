-- reverse: modify "banners" table
ALTER TABLE "banners" DROP CONSTRAINT "fk_banners_country", DROP COLUMN "country_code";
