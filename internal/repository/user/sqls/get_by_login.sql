SELECT id, login, password_hash, display_name
FROM users
WHERE login = $1
