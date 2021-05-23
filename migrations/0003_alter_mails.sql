ALTER TABLE mails ADD deleted_by_sender BOOLEAN DEFAULT FALSE;
ALTER TABLE mails ADD deleted_by_recipient BOOLEAN DEFAULT FALSE;