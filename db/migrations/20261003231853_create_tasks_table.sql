CREATE TABLE tasks (
   id          UUID PRIMARY KEY,
   family_id   UUID NOT NULL REFERENCES families (id),
   title       TEXT NOT NULL,
   description TEXT NOT NULL DEFAULT '',
   points      BIGINT NOT NULL CHECK (points > 0),
   created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
   updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_tasks_family_id
    ON tasks (family_id);
