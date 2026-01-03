SELECT s.session_id, s.token, s.staff_id, s.created_at, s.expires_at,
       st.id, st.username, st.email, st.first_name, st.last_name, st.role, st.is_active
FROM sessions s
JOIN staff st ON s.staff_id = st.id
WHERE s.session_id = $1
