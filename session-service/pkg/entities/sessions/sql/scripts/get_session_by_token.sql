SELECT session_id, token
FROM sessions
WHERE token = $1
