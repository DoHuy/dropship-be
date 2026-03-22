-- reverse: modify "variants" table
ALTER TABLE "variants" ADD COLUMN "weight_gram" bigint NULL DEFAULT 0;
