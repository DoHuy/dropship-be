-- modify "sliders" table
ALTER TABLE "sliders" ADD COLUMN "country_code" character(2) NOT NULL, ADD CONSTRAINT "fk_sliders_country" FOREIGN KEY ("country_code") REFERENCES "countries" ("code") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- create index "idx_sliders_country_code" to table: "sliders"
CREATE INDEX "idx_sliders_country_code" ON "sliders" ("country_code");
