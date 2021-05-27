ALTER TABLE mails ADD deleted_by_sender BOOLEAN DEFAULT FALSE;
ALTER TABLE mails ADD deleted_by_recipient BOOLEAN DEFAULT FALSE;
ALTER TABLE dialogues ADD CONSTRAINT dialogues_owner_fkey FOREIGN KEY owner REFERENCES users (username) ON DELETE CASCADE;
ALTER TABLE dialogues ADD body TEXT;