-- create "sliders" table
CREATE TABLE "sliders" (
  "id" bigserial NOT NULL,
  "title" character varying(255) NOT NULL,
  "image_url" character varying(500) NOT NULL,
  "sub_text" text NULL,
  "description" text NULL,
  "position" bigint NULL DEFAULT 0,
  "is_active" boolean NULL DEFAULT true,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- create index "idx_sliders_deleted_at" to table: "sliders"
CREATE INDEX "idx_sliders_deleted_at" ON "sliders" ("deleted_at");
