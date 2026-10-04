CREATE TYPE task_participation_mode AS ENUM (
    'each_child',
    'shared'
    );

CREATE TABLE season_tasks
(
    id                 UUID PRIMARY KEY,
    season_id          UUID                    NOT NULL REFERENCES seasons (id),
    task_id            UUID                    NOT NULL REFERENCES tasks (id),
    title              TEXT                    NOT NULL,
    description        TEXT                    NOT NULL DEFAULT '',
    points             BIGINT                  NOT NULL CHECK (points > 0),
    participation_mode task_participation_mode NOT NULL,
    created_at         TIMESTAMPTZ             NOT NULL DEFAULT now(),
    updated_at         TIMESTAMPTZ             NOT NULL DEFAULT now(),

    CONSTRAINT uq_season_tasks_season_task
        UNIQUE (season_id, task_id)
);
