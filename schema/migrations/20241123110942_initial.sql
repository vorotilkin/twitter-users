-- Create "user" table
CREATE TABLE "user" ("id" serial NOT NULL, "name" text NOT NULL, "password_hash" text NOT NULL, "username" text NOT NULL, "email" text NOT NULL, "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP, "updated_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY ("id"));
-- Create index "idx_id" to table: "user"
CREATE INDEX "idx_id" ON "user" ("id");
