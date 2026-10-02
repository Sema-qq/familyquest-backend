CREATE TYPE family_member_role AS ENUM ('parent', 'child');

CREATE TABLE family_members (
    id         UUID PRIMARY KEY,
    family_id  UUID NOT NULL REFERENCES families (id),
    user_id    UUID NOT NULL UNIQUE REFERENCES users (id),
    role       family_member_role NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);