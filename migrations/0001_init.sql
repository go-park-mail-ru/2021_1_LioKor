CREATE EXTENSION IF NOT EXISTS CITEXT;

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username CITEXT UNIQUE NOT NULL,
    password_hash CITEXT NOT NULL,
    avatar_url VARCHAR(128),
    fullname VARCHAR(128),
    reserve_email CITEXT
);

CREATE TABLE IF NOT EXISTS sessions (
    user_id INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    token CITEXT PRIMARY KEY,
    expiration TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS mails (
    id BIGSERIAL PRIMARY KEY,
    sender CITEXT NOT NULL,
    recipient CITEXT NOT NULL,
    subject TEXT,
    body TEXT,
    received_date  TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    unread BOOLEAN DEFAULT TRUE,
    status INT DEFAULT 1
);

CREATE TABLE IF NOT EXISTS folders (
    id BIGSERIAL PRIMARY KEY,
    folder_name CITEXT NOT NULL,
    owner INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    UNIQUE(folder_name, owner)
);

CREATE TABLE IF NOT EXISTS dialogues (
    id BIGSERIAL PRIMARY KEY,
    owner CITEXT NOT NULL,
    other CITEXT NOT NULL,
    last_mail_id INTEGER REFERENCES mails (id) ON DELETE SET NULL,
    received_date  TIMESTAMP,
    unread INT DEFAULT 0,
    folder INTEGER REFERENCES folders (id) DEFAULT NULL,
    UNIQUE (owner, other)
);


CREATE OR REPLACE FUNCTION add_dialogue()
    RETURNS TRIGGER
    AS $add_dialogue$
DECLARE
BEGIN
    IF EXISTS (
        SELECT * FROM dialogues WHERE (owner=NEW.sender AND other=NEW.recipient) OR (other=NEW.sender AND owner=NEW.recipient) LIMIT 1
    ) THEN
        UPDATE dialogues SET last_mail_id=NEW.id, received_date=NEW.received_date WHERE owner=NEW.sender AND other=NEW.recipient;
        UPDATE dialogues SET last_mail_id=NEW.id, received_date=NEW.received_date, unread=(unread + 1) WHERE other=NEW.sender AND owner=NEW.recipient;
    ELSE
        INSERT INTO dialogues(owner, other, last_mail_id, received_date) values (NEW.sender, NEW.recipient, NEW.id, NEW.received_date);
        INSERT INTO dialogues(owner, other, last_mail_id, received_date, unread) values (NEW.recipient, NEW.sender, NEW.id, NEW.received_date, 1);
    END IF;
RETURN NEW;
END;
$add_dialogue$ LANGUAGE plpgsql;

CREATE TRIGGER add_dialogue
    AFTER INSERT
    ON mails
    FOR EACH ROW
    EXECUTE PROCEDURE add_dialogue();
