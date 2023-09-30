-- migrate:up
CREATE TABLE public."person" (
    "id" BIGSERIAL,
    "name" VARCHAR NOT NULL,
    "created_at" TIMESTAMPTZ DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT "person_pkey" PRIMARY KEY ("id")
);

-- migrate:down
DROP TABLE IF EXISTS public."person";
