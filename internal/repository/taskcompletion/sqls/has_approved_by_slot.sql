SELECT COUNT(*) AS count
FROM task_completions
WHERE slot_id = $1
  AND status = 'approved'
