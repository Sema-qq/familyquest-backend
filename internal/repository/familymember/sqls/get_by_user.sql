SELECT
    fm.id,
    fm.family_id,
    fm.user_id,
    fm.role
FROM family_members fm
WHERE fm.user_id = $1
