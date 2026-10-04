WITH updated AS (
    UPDATE task_completions
    SET status = 'approved',
        reviewed_by = $2,
        reviewed_at = $3,
        updated_at = now()
    WHERE id = $1
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
SELECT u.id,
       u.member_id,
       u.slot_id,
       sts.season_task_id,
       u.performed_on,
       u.status,
       u.child_comment,
       u.parent_comment,
       u.reviewed_by,
       u.reviewed_at
FROM updated u
JOIN season_task_slots sts ON sts.id = u.slot_id
