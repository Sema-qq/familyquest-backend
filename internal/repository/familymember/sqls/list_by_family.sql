SELECT
    fm.id,
    fm.family_id,
    u.id AS user_id,
    u.login,
    u.display_name,
    fm.role
FROM family_members fm
JOIN users u ON u.id = fm.user_id
WHERE fm.family_id = $1
ORDER BY fm.created_at
