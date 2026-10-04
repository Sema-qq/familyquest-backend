SELECT sts.id,
       sts.season_task_id,
       sts.position,
       sts.available_from,
       sts.available_until
FROM season_task_slots sts
JOIN season_tasks st ON st.id = sts.season_task_id
WHERE st.season_id = $1
ORDER BY st.created_at, sts.position
