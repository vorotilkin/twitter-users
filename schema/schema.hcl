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
table "follow" {
  schema = schema.public

  column "user_id" {
    null = false
    type = integer
    comment = "Идентификатор пользователя, который подписывается"
  }

  column "following_user_id" {
    null = false
    type = integer
    comment = "Идентификатор пользователя, на которого подписываются"
  }

  column "created_at" {
    null    = false
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
    comment = "Дата и время создания подписки"
  }

  primary_key {
    columns = [column.user_id, column.following_user_id]
  }

  foreign_key "fk_follow_user_id" {
    columns    = [column.user_id]
    ref_columns = [table.user.column.id]
    on_delete = CASCADE
  }

  foreign_key "fk_follow_following_user_id" {
    columns    = [column.following_user_id]
    ref_columns = [table.user.column.id]
    on_delete = CASCADE
  }

  index "idx_follow_user_id" {
    columns = [column.user_id]
  }

  index "idx_follow_following_user_id" {
    columns = [column.following_user_id]
  }
}
schema "public" {
  comment = "standard public schema"
}
