SELECT f.id, f.name, f.timezone, f.owner_id, fm.role
FROM families f
JOIN family_members fm ON fm.family_id = f.id
WHERE fm.user_id = $1
