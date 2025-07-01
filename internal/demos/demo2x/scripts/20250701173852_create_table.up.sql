CREATE TABLE "users" ("id" bigserial,"created_at" timestamptz,"updated_at" timestamptz,"deleted_at" timestamptz,"username" text,"nickname" text,"rank" bigint,"score" decimal,PRIMARY KEY ("id"),CONSTRAINT "uni_users_username" UNIQUE ("username"));

CREATE INDEX IF NOT EXISTS "idx_users_deleted_at" ON "users" ("deleted_at");

CREATE TABLE "infos" ("id" bigserial,"created_at" timestamptz,"updated_at" timestamptz,"deleted_at" timestamptz,"name" text,"cate" bigint,PRIMARY KEY ("id"),CONSTRAINT "uni_infos_name" UNIQUE ("name"));

CREATE INDEX IF NOT EXISTS "idx_infos_deleted_at" ON "infos" ("deleted_at");
