SELECT tc.id,
       tc.member_id,
       tc.slot_id,
       sts.season_task_id,
       tc.performed_on,
       tc.status,
       tc.child_comment,
       tc.parent_comment,
       tc.reviewed_by,
       tc.reviewed_at
FROM task_completions tc
JOIN season_task_slots sts ON sts.id = tc.slot_id
WHERE tc.member_id = $1
  AND tc.slot_id = $2
  AND tc.status = 'rejected'
