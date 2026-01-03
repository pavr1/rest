UPDATE sessions SET token = $2, expires_at = $3 WHERE session_id = $1
