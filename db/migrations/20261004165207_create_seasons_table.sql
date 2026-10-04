CREATE TYPE season_status AS ENUM (
    'draft',
    'active',
    'completed'
    );

CREATE TABLE seasons (
    id           UUID PRIMARY KEY,
    family_id    UUID NOT NULL REFERENCES families (id),
    title        TEXT NOT NULL,
    starts_at    DATE NOT NULL,
    ends_at      DATE NOT NULL,
    timezone     TEXT NOT NULL,
    status       season_status NOT NULL DEFAULT 'draft',
    completed_at TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT chk_seasons_dates CHECK (ends_at >= starts_at)
);

CREATE INDEX idx_seasons_family_id
    ON seasons (family_id);

CREATE UNIQUE INDEX uq_seasons_one_active_per_family
    ON seasons (family_id)
    WHERE status = 'active';
