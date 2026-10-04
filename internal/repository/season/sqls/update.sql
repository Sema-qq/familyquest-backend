UPDATE seasons
SET status = COALESCE($2::season_status, status),
    completed_at = COALESCE($3::timestamptz, completed_at),
    updated_at = now()
WHERE id = $1
