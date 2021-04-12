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
    received_date  TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS dialogues (
    user_1 CITEXT NOT NULL,
    user_2 CITEXT NOT NULL,
    mail_id INTEGER REFERENCES mails (id) ON DELETE SET NULL,
    received_date  TIMESTAMP,
    UNIQUE (user_1, user_2)
);

CREATE OR REPLACE FUNCTION add_dialogue()
    RETURNS TRIGGER
    AS $add_dialogue$
DECLARE
    bigger CITEXT;
    smaller CITEXT;
BEGIN
    IF NEW.sender > NEW.recipient THEN
        bigger := NEW.sender;
        smaller := NEW.recipient;
    ELSE
        bigger := NEW.recipient;
        smaller := NEW.sender;
    END IF;
    IF EXISTS (
        SELECT * FROM dialogues WHERE user_1=bigger AND user_2=smaller LIMIT 1
    ) THEN
        UPDATE dialogues SET mail_id=NEW.id, received_date=NEW.received_date WHERE user_1=bigger AND user_2=smaller;
    ELSE
        INSERT INTO dialogues(user_1, user_2, mail_id, received_date) values (bigger, smaller, NEW.id, NEW.received_date);
    END IF;
RETURN NEW;
END;
$add_dialogue$ LANGUAGE plpgsql;

CREATE TRIGGER add_dialogue
    AFTER INSERT
    ON mails
    FOR EACH ROW
    EXECUTE PROCEDURE add_dialogue();