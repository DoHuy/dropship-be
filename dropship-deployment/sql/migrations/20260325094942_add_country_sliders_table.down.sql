-- reverse: create index "idx_sliders_country_code" to table: "sliders"
DROP INDEX "idx_sliders_country_code";
-- reverse: modify "sliders" table
ALTER TABLE "sliders" DROP CONSTRAINT "fk_sliders_country", DROP COLUMN "country_code";
