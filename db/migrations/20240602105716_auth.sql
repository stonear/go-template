-- migrate:up
CREATE TABLE public."auth" (
    "id" UUID PRIMARY KEY DEFAULT GEN_RANDOM_UUID(),
    "username" VARCHAR NOT NULL,
    "password" VARCHAR NOT NULL,
    "created_at" TIMESTAMPTZ DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ DEFAULT NOW()
);
-- admin:admin123 (bcrypt)
-- example only, do not use in production
INSERT INTO public."auth"("username", "password") 
VALUES ('admin', '$2a$12$b7zm4mQ0JNu3Pr4FlM3cuudcf8SfpVrHAT/2MCZcF6yg3sqewPdC2');

-- migrate:down
DROP TABLE IF EXISTS public."auth";