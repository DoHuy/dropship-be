-- modify "banners" table
ALTER TABLE "banners" ADD COLUMN "video_url" character varying(500) NULL, ADD COLUMN "alt" character varying(255) NULL, ADD COLUMN "description" text NULL, ADD COLUMN "video_type" character varying(50) NULL;
