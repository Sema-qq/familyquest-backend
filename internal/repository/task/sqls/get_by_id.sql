SELECT id,
       family_id,
       title,
       description,
       points
FROM tasks
WHERE id = $1
