CREATE TYPE zapiski.note_source AS ENUM ('manual', 'message');

CREATE TABLE zapiski.notes (
    id                bigserial PRIMARY KEY,
    version bigint NOT NULL DEFAULT 1,
    user_id           bigint NOT NULL,
    chat_id           bigint REFERENCES zapiski.chats(id)    ON DELETE SET NULL,
    source_message_id bigint REFERENCES zapiski.messages(id) ON DELETE SET NULL,
    source            zapiski.note_source NOT NULL DEFAULT 'manual',        -- manual | message
    title             text NOT NULL,
    content           text NOT NULL,
    created_at        timestamptz NOT NULL DEFAULT now(),
    updated_at        timestamptz,
    deleted_at        timestamptz,

    CONSTRAINT notes_source_check
        CHECK (source IN ('manual', 'message')),
    CONSTRAINT notes_message_source_check
        CHECK (source <> 'message' OR source_message_id IS NOT NULL)
);

CREATE UNIQUE INDEX notes_user_source_uniq
    ON zapiski.notes (user_id, source_message_id)
    WHERE source_message_id IS NOT NULL AND deleted_at IS NULL;

CREATE INDEX notes_user_created_idx
    ON zapiski.notes (user_id, created_at DESC)
    WHERE deleted_at IS NULL;