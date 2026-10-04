SELECT id,
       season_id,
       task_id,
       title,
       description,
       points,
       participation_mode
FROM season_tasks
WHERE season_id = $1
ORDER BY created_at, id
