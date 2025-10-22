-- +goose Up
-- +goose StatementBegin
-- Step 1: Drop foreign key constraints that depend on Player.uid
ALTER TABLE Session DROP CONSTRAINT session_uid_fkey;
ALTER TABLE Wallet DROP CONSTRAINT wallet_uid_fkey;

-- Step 2: Drop primary key constraint on Player.uid
ALTER TABLE Player DROP CONSTRAINT player_pkey;

-- Step 3: Alter Player.uid to VARCHAR(36)
ALTER TABLE Player ALTER COLUMN uid TYPE VARCHAR(36);

-- Step 4: Re-add primary key constraint
ALTER TABLE Player ADD PRIMARY KEY (uid);

-- Step 5: Alter Session.uid to VARCHAR(36)
ALTER TABLE Session ALTER COLUMN uid TYPE VARCHAR(36);

-- Step 6: Re-add foreign key constraint to Session.uid
ALTER TABLE Session ADD CONSTRAINT session_uid_fkey FOREIGN KEY (uid) REFERENCES Player(uid) ON DELETE CASCADE;

-- Step 7: Alter Wallet.uid to VARCHAR(36)
ALTER TABLE Wallet ALTER COLUMN uid TYPE VARCHAR(36);

-- Step 8: Re-add foreign key constraint to Wallet.uid
ALTER TABLE Wallet ADD CONSTRAINT wallet_uid_fkey FOREIGN KEY (uid) REFERENCES Player(uid) ON DELETE CASCADE;

-- +goose StatementEnd
