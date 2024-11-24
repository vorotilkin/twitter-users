-- Modify "user" table
ALTER TABLE "user" ADD COLUMN "bio" text NULL, ADD COLUMN "email_verified" timestamp NULL, ADD COLUMN "image" text NULL, ADD COLUMN "cover_image" text NULL, ADD COLUMN "profile_image" text NULL, ADD COLUMN "has_notification" boolean NOT NULL DEFAULT false;
