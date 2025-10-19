-- +goose Up
-- +goose StatementBegin

-- Step 1: Drop foreign key constraint from SessionState
ALTER TABLE SessionState DROP CONSTRAINT SessionState_session_id_fkey;

-- Step 2: Drop primary key constraint from Session
ALTER TABLE Session DROP CONSTRAINT Session_pkey;

-- Step 3: Alter session_id column in Session to VARCHAR(36)
ALTER TABLE Session ALTER COLUMN session_id TYPE VARCHAR(36);

-- Step 4: Recreate primary key on Session
ALTER TABLE Session ADD PRIMARY KEY (session_id);

-- Step 5: Alter session_id column in SessionState to VARCHAR(36)
ALTER TABLE SessionState ALTER COLUMN session_id TYPE VARCHAR(36);

-- Step 6: Recreate foreign key constraint in SessionState
ALTER TABLE SessionState ADD CONSTRAINT SessionState_session_id_fkey
    FOREIGN KEY (session_id) REFERENCES Session(session_id) ON DELETE CASCADE;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Step 1: Drop foreign key constraint from SessionState
ALTER TABLE SessionState DROP CONSTRAINT SessionState_session_id_fkey;

-- Step 2: Drop primary key constraint from Session
ALTER TABLE Session DROP CONSTRAINT Session_pkey;

-- Step 3: Alter session_id column in Session back to SERIAL (INTEGER)
ALTER TABLE Session ALTER COLUMN session_id TYPE INTEGER USING session_id::INTEGER;

-- Step 4: Recreate primary key on Session
ALTER TABLE Session ADD PRIMARY KEY (session_id);

-- Step 5: Alter session_id column in SessionState back to INTEGER
ALTER TABLE SessionState ALTER COLUMN session_id TYPE INTEGER USING session_id::INTEGER;

-- Step 6: Recreate foreign key constraint in SessionState
ALTER TABLE SessionState ADD CONSTRAINT SessionState_session_id_fkey
    FOREIGN KEY (session_id) REFERENCES Session(session_id) ON DELETE CASCADE;

-- +goose StatementEnd
