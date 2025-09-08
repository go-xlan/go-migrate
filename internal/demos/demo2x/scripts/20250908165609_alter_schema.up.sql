ALTER TABLE "users" ALTER COLUMN "rank" TYPE text USING "rank"::text;

ALTER TABLE "infos" ALTER COLUMN "cate" TYPE bigint USING "cate"::bigint;
