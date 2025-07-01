ALTER TABLE "users" ALTER COLUMN "score" TYPE text USING "score"::text;

ALTER TABLE "infos" ALTER COLUMN "cate" TYPE text USING "cate"::text;
