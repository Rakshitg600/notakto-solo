-- +goose Up
-- +goose StatementBegin
-- 1) Drop current_player and expired attributes from Session table
ALTER TABLE Session
    DROP COLUMN IF EXISTS current_player,
    DROP COLUMN IF EXISTS expired;

-- 2) Drop game_history from SessionState table
ALTER TABLE SessionState
    DROP COLUMN IF EXISTS game_history;

-- 3) Move gameover, winner, board_size, number_of_boards, difficulty from SessionState to Session
-- First, add these columns to Session
ALTER TABLE Session
    ADD COLUMN IF NOT EXISTS gameover BOOLEAN DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS winner BOOLEAN,
    ADD COLUMN IF NOT EXISTS board_size INTEGER,
    ADD COLUMN IF NOT EXISTS number_of_boards INTEGER,
    ADD COLUMN IF NOT EXISTS difficulty INTEGER;

-- Then, drop these columns from SessionState
ALTER TABLE SessionState
    DROP COLUMN IF EXISTS gameover,
    DROP COLUMN IF EXISTS winner,
    DROP COLUMN IF EXISTS board_size,
    DROP COLUMN IF EXISTS number_of_boards,
    DROP COLUMN IF EXISTS difficulty;

-- 4) Drop boards column and add it back as INTEGER[] (1D array of integers)
ALTER TABLE SessionState
    DROP COLUMN IF EXISTS boards,
    ADD COLUMN boards INTEGER[];

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Revert changes in reverse order
-- 4) Revert boards: drop INTEGER[] column and restore as TEXT[][]
ALTER TABLE SessionState
    DROP COLUMN IF EXISTS boards,
    ADD COLUMN boards TEXT[][];

-- 3) Move gameover, winner, board_size, number_of_boards, difficulty back to SessionState
-- Add columns back to SessionState
ALTER TABLE SessionState
    ADD COLUMN IF NOT EXISTS gameover BOOLEAN DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS winner VARCHAR(50),
    ADD COLUMN IF NOT EXISTS board_size INTEGER,
    ADD COLUMN IF NOT EXISTS number_of_boards INTEGER,
    ADD COLUMN IF NOT EXISTS difficulty INTEGER;

-- Drop columns from Session
ALTER TABLE Session
    DROP COLUMN IF EXISTS gameover,
    DROP COLUMN IF EXISTS winner,
    DROP COLUMN IF EXISTS board_size,
    DROP COLUMN IF EXISTS number_of_boards,
    DROP COLUMN IF EXISTS difficulty;

-- 2) Restore game_history in SessionState
ALTER TABLE SessionState
    ADD COLUMN IF NOT EXISTS game_history TEXT[][][];

-- 1) Restore current_player and expired in Session
ALTER TABLE Session
    ADD COLUMN IF NOT EXISTS current_player INTEGER,
    ADD COLUMN IF NOT EXISTS expired BOOLEAN DEFAULT FALSE;

-- +goose StatementEnd