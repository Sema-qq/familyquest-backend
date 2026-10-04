WITH inserted AS (
    INSERT INTO task_completions (id, member_id, slot_id, performed_on, child_comment)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING id,
              member_id,
              slot_id,
              performed_on,
              status,
              child_comment,
              parent_comment,
              reviewed_by,
              reviewed_at
)
SELECT i.id,
       i.member_id,
       i.slot_id,
       sts.season_task_id,
       i.performed_on,
       i.status,
       i.child_comment,
       i.parent_comment,
       i.reviewed_by,
       i.reviewed_at
FROM inserted i
JOIN season_task_slots sts ON sts.id = i.slot_id
