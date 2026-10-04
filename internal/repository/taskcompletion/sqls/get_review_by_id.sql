SELECT tc.id,
       tc.member_id,
       tc.slot_id,
       sts.season_task_id,
       tc.performed_on,
       tc.status,
       tc.child_comment,
       tc.parent_comment,
       tc.reviewed_by,
       tc.reviewed_at,
       st.season_id,
       s.family_id AS season_family_id,
       s.status AS season_status,
       st.participation_mode
FROM task_completions tc
JOIN season_task_slots sts ON sts.id = tc.slot_id
JOIN season_tasks st ON st.id = sts.season_task_id
JOIN seasons s ON s.id = st.season_id
WHERE tc.id = $1
