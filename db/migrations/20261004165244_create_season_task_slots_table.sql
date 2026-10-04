CREATE TABLE season_task_slots
(
    id              UUID PRIMARY KEY,
    season_task_id  UUID        NOT NULL REFERENCES season_tasks (id),
    position        INTEGER     NOT NULL CHECK (position > 0),
    available_from  DATE        NOT NULL,
    available_until DATE        NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT chk_season_task_slots_dates CHECK (available_until >= available_from),

    CONSTRAINT uq_season_task_slots_position UNIQUE (season_task_id, position)
);