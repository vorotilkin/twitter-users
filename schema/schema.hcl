table "user" {
  schema = schema.public
  column "id" {
    null = false
    type = serial
  }
  column "name" {
    null = false
    type = text
  }
  column "password_hash" {
    null = false
    type = text
  }
  column "username" {
    null = false
    type = text
  }
  column "email" {
    null = false
    type = text
  }
  column "created_at" {
    null    = false
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = false
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  column "bio" {
    null = true
    type = text
  }
  column "email_verified" {
    null = true
    type = timestamp
  }
  column "image" {
    null = true
    type = text
  }
  column "cover_image" {
    null = true
    type = text
  }
  column "profile_image" {
    null = true
    type = text
  }
  column "has_notification" {
    null    = false
    type    = boolean
    default = false
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_id" {
    columns = [column.id]
  }
  unique "user_pk" {
    columns = [column.email]
  }
}
schema "public" {
  comment = "standard public schema"
}
