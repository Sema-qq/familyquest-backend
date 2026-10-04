SELECT tc.id,
       tc.member_id,
       u.display_name,
       tc.slot_id,
       st.id AS season_task_id,
       st.title AS task_title,
       tc.performed_on,
       tc.status,
       st.points,
       tc.child_comment,
       tc.parent_comment,
       tc.reviewed_by,
       tc.reviewed_at
FROM task_completions tc
JOIN family_members fm ON fm.id = tc.member_id
JOIN users u ON u.id = fm.user_id
JOIN season_task_slots sts ON sts.id = tc.slot_id
JOIN season_tasks st ON st.id = sts.season_task_id
WHERE st.season_id = $1
  AND ($2::uuid IS NULL OR tc.member_id = $2)
  AND ($3::task_completion_status IS NULL OR tc.status = $3)
ORDER BY tc.created_at DESC
