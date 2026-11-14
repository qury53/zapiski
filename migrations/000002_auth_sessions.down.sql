DROP INDEX IF EXISTS zapiski.zapiski_user_sessions_parent_session_id_idx;
DROP INDEX IF EXISTS zapiski.zapiski_user_sessions_ip_address_idx;
DROP INDEX IF EXISTS zapiski.zapiski_user_sessions_expires_at_idx;
DROP INDEX IF EXISTS zapiski.zapiski_user_sessions_user_id_active_idx;
DROP INDEX IF EXISTS zapiski.zapiski_user_sessions_user_id_idx;
DROP INDEX IF EXISTS zapiski.zapiski_user_sessions_jti_uidx;

DROP TABLE IF EXISTS zapiski.user_sessions;