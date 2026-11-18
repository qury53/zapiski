CREATE TABLE zapiski.user_sessions (
    id BIGSERIAL PRIMARY KEY,

    user_id BIGINT NOT NULL
        REFERENCES zapiski.users(id)
        ON DELETE CASCADE,

    jti UUID NOT NULL,
    refresh_token_hash VARCHAR(255) NOT NULL,

    user_agent TEXT,
    ip_address INET,

    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ,
    revoked_reason VARCHAR(64),

    parent_session_id BIGINT
        REFERENCES zapiski.user_sessions(id)
        ON DELETE SET NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CHECK (revoked_reason IS NULL OR revoked_at IS NOT NULL)
);

CREATE UNIQUE INDEX zapiski_user_sessions_jti_uidx
    ON zapiski.user_sessions (jti);

CREATE INDEX zapiski_user_sessions_user_id_idx
    ON zapiski.user_sessions (user_id);

CREATE INDEX zapiski_user_sessions_user_id_active_idx
    ON zapiski.user_sessions (user_id, expires_at)
    WHERE revoked_at IS NULL;

CREATE INDEX zapiski_user_sessions_expires_at_idx
    ON zapiski.user_sessions (expires_at);

CREATE INDEX zapiski_user_sessions_ip_address_idx
    ON zapiski.user_sessions (ip_address)
    WHERE ip_address IS NOT NULL;

CREATE INDEX zapiski_user_sessions_parent_session_id_idx
    ON zapiski.user_sessions (parent_session_id)
    WHERE parent_session_id IS NOT NULL;