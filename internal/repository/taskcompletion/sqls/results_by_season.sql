SELECT fm.id AS member_id,
       u.display_name,
       COUNT(st.id) FILTER (WHERE tc.status = 'approved') AS approved_count,
       COALESCE(SUM(st.points) FILTER (WHERE tc.status = 'approved'), 0) AS points
FROM family_members fm
JOIN users u ON u.id = fm.user_id
LEFT JOIN task_completions tc ON tc.member_id = fm.id
LEFT JOIN season_task_slots sts ON sts.id = tc.slot_id
LEFT JOIN season_tasks st ON st.id = sts.season_task_id
                         AND st.season_id = $1
WHERE fm.family_id = (
    SELECT family_id
    FROM seasons
    WHERE id = $1
)
  AND fm.role = 'child'
  AND ($2::uuid IS NULL OR fm.id = $2)
GROUP BY fm.id, u.display_name
ORDER BY u.display_name
