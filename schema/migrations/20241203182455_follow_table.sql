-- Create "follow" table
CREATE TABLE "follow" ("user_id" integer NOT NULL, "following_user_id" integer NOT NULL, "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY ("user_id", "following_user_id"), CONSTRAINT "fk_follow_following_user_id" FOREIGN KEY ("following_user_id") REFERENCES "user" ("id") ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT "fk_follow_user_id" FOREIGN KEY ("user_id") REFERENCES "user" ("id") ON UPDATE NO ACTION ON DELETE CASCADE);
-- Create index "idx_follow_following_user_id" to table: "follow"
CREATE INDEX "idx_follow_following_user_id" ON "follow" ("following_user_id");
-- Create index "idx_follow_user_id" to table: "follow"
CREATE INDEX "idx_follow_user_id" ON "follow" ("user_id");
-- Set comment to column: "user_id" on table: "follow"
COMMENT ON COLUMN "follow"."user_id" IS 'Идентификатор пользователя, который подписывается';
-- Set comment to column: "following_user_id" on table: "follow"
COMMENT ON COLUMN "follow"."following_user_id" IS 'Идентификатор пользователя, на которого подписываются';
-- Set comment to column: "created_at" on table: "follow"
COMMENT ON COLUMN "follow"."created_at" IS 'Дата и время создания подписки';
