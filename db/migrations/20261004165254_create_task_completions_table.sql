CREATE TYPE task_completion_status AS ENUM (
    'submitted',
    'approved',
    'rejected'
    );

CREATE TABLE task_completions
(
    id             UUID PRIMARY KEY,
    member_id      UUID                   NOT NULL REFERENCES family_members (id),
    slot_id        UUID                   NOT NULL REFERENCES season_task_slots (id),
    performed_on   DATE                   NOT NULL,
    status         task_completion_status NOT NULL DEFAULT 'submitted',
    child_comment  TEXT                   NOT NULL DEFAULT '',
    parent_comment TEXT                   NOT NULL DEFAULT '',
    reviewed_by    UUID REFERENCES family_members (id),
    reviewed_at    TIMESTAMPTZ,
    created_at     TIMESTAMPTZ            NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ            NOT NULL DEFAULT now(),

    CONSTRAINT uq_task_completions_member_slot
        UNIQUE (member_id, slot_id)
);

CREATE INDEX idx_task_completions_slot_status
    ON task_completions (slot_id, status);
