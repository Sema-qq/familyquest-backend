SELECT COUNT(*) AS count
FROM task_completions tc
JOIN season_task_slots sts ON sts.id = tc.slot_id
JOIN season_tasks st ON st.id = sts.season_task_id
WHERE st.season_id = $1
  AND tc.status = 'submitted'
