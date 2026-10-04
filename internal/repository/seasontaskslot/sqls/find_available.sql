SELECT sts.id,
       sts.season_task_id,
       sts.position,
       sts.available_from,
       sts.available_until
FROM season_task_slots sts
WHERE sts.season_task_id = $1
  AND sts.available_from <= $3
  AND sts.available_until >= $3
  AND NOT EXISTS (
      SELECT 1
      FROM task_completions tc
      WHERE tc.slot_id = sts.id
        AND tc.member_id = $2
        AND tc.status IN ('submitted', 'approved')
  )
  AND (
      $4 = FALSE
      OR NOT EXISTS (
          SELECT 1
          FROM task_completions tc
          WHERE tc.slot_id = sts.id
            AND tc.status IN ('submitted', 'approved')
      )
  )
ORDER BY sts.position
LIMIT 1
