CREATE TABLE families (
    id         UUID PRIMARY KEY,
    name       TEXT NOT NULL,
    timezone   TEXT NOT NULL,
    owner_id   UUID NOT NULL UNIQUE REFERENCES users (id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);