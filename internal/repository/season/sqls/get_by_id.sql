SELECT id,
       family_id,
       title,
       starts_at,
       ends_at,
       timezone,
       status,
       completed_at
FROM seasons
WHERE id = $1
