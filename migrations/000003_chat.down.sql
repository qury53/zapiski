DROP TRIGGER IF EXISTS messages_touch_chat ON zapiski.messages;

DROP FUNCTION IF EXISTS zapiski.touch_chat();

DROP TABLE IF EXISTS zapiski.messages;

DROP TABLE IF EXISTS zapiski.chats;

DROP TYPE IF EXISTS zapiski.msg_role;
DROP TYPE IF EXISTS zapiski.chat_mode;