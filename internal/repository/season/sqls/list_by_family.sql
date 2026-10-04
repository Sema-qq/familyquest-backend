SELECT id,
       family_id,
       title,
       starts_at,
       ends_at,
       timezone,
       status,
       completed_at
FROM seasons
WHERE family_id = $1
ORDER BY starts_at DESC, created_at DESC
