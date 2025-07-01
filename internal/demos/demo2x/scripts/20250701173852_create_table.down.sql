-- reverse -- CREATE INDEX IF NOT EXISTS "idx_infos_deleted_at" ON "infos" ("deleted_at");
SELECT TODO / PANIC / RAISE / THROW; -- DROP INDEX; -- TODO

-- reverse -- CREATE TABLE "infos" ("id" bigserial,"created_at" timestamptz,"updated_at" timestamptz,"deleted_at" timestamptz,"name" text,"cate" bigint,PRIMARY KEY ("id"),CONSTRAINT "uni_infos_name" UNIQUE ("name"));
SELECT TODO / PANIC / RAISE / THROW; -- DROP TABLE; -- TODO

-- reverse -- CREATE INDEX IF NOT EXISTS "idx_users_deleted_at" ON "users" ("deleted_at");
SELECT TODO / PANIC / RAISE / THROW; -- DROP INDEX; -- TODO

-- reverse -- CREATE TABLE "users" ("id" bigserial,"created_at" timestamptz,"updated_at" timestamptz,"deleted_at" timestamptz,"username" text,"nickname" text,"rank" bigint,"score" decimal,PRIMARY KEY ("id"),CONSTRAINT "uni_users_username" UNIQUE ("username"));
SELECT TODO / PANIC / RAISE / THROW; -- DROP TABLE; -- TODO
