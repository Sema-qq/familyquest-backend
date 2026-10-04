SELECT COUNT(*) AS count
FROM season_task_slots sts
JOIN season_tasks st ON st.id = sts.season_task_id
WHERE st.season_id = $1
