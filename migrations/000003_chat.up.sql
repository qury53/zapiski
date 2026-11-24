CREATE TYPE zapiski.chat_mode AS ENUM ('support', 'full');
CREATE TYPE zapiski.msg_role AS ENUM ('unspecified', 'system', 'user', 'assistant', 'tool')

CREATE FUNCTION zapiski.touch_chat() RETURNS trigger AS $$
BEGIN
    UPDATE zapiski.chats SET updated_at = now() WHERE id = NEW.chat_id;
    RETURN NEW;
END $$ LANGUAGE plpgsql;

CREATE TRIGGER messages_touch_chat
AFTER INSERT OR UPDATE OR DELETE ON zapiski.messages
FOR EACH ROW EXECUTE FUNCTION zapiski.touch_chat();

CREATE TABLE zapiski.chats (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL 
        REFERENCES zapiski.users(id)
        ON DELETE CASCADE,
    mode zapiski.chat_mode NOT NULL DEFAULT 'full',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
)

CREATE TABLE zapiski.messages (
    id BIGSERIAL PRIMARY KEY,
    chat_id BIGINT NOT NULL
        REFERENCES zapiski.chats(id)
        ON DELETE CASCADE,
    role zapiski.msg_role NOT NULL,
    content TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
)

CREATE INDEX ON zapiski.chats (user_id, updated_at DESC)
    WHERE deleted_at IS NULL;

CREATE INDEX ON zapiski.messages (chat_id, id);