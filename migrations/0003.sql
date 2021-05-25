ALTER TABLE mails
    ADD deleted_by_sender BOOLEAN DEFAULT false,
    ADD deleted_by_recipient BOOLEAN DEFAULT false;
