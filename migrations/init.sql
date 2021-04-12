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
    user_id INTEGER NOT NULL REFERENCES users (id),
    token CITEXT PRIMARY KEY,
    expiration TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS mails (
    sender CITEXT NOT NULL,
    recipient CITEXT NOT NULL,
    subject TEXT,
    body TEXT,
    received_date  TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
