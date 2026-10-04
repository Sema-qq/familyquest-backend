SELECT id,
       family_id,
       title,
       description,
       points
FROM tasks
WHERE family_id = $1
ORDER BY created_at DESC
