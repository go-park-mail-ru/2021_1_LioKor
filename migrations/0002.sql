ALTER TABLE dialogues
    ALTER COLUMN received_date TYPE TIMESTAMP WITH TIME ZONE;
ALTER TABLE dialogues
    ALTER COLUMN received_date SET DEFAULT NOW();