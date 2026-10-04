SELECT st.id,
       st.season_id,
       st.task_id,
       st.title,
       st.description,
       st.points,
       st.participation_mode,
       s.family_id AS season_family_id,
       s.starts_at AS season_starts_at,
       s.ends_at AS season_ends_at,
       s.timezone AS season_timezone,
       s.status AS season_status
FROM season_tasks st
JOIN seasons s ON s.id = st.season_id
WHERE st.id = $1
