-- migrate:up
CREATE TABLE public."person" (
    "id" UUID PRIMARY KEY DEFAULT GEN_RANDOM_UUID(),
    "name" VARCHAR NOT NULL,
    "created_at" TIMESTAMPTZ DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ DEFAULT NOW()
);

-- migrate:down
DROP TABLE IF EXISTS public."person";
