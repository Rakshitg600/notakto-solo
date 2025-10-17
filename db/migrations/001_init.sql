-- +goose Up
-- +goose StatementBegin
-- Create Player table
CREATE TABLE Player (
    uid SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    profile_pic TEXT
);

-- Create Session table
CREATE TABLE Session (
    session_id SERIAL PRIMARY KEY,
    uid INTEGER NOT NULL,
    expired BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (uid) REFERENCES Player(uid) ON DELETE CASCADE
);

-- Create SessionState table
CREATE TABLE SessionState (
    session_id INTEGER PRIMARY KEY,
    boards TEXT[][], -- 2D array of characters
    current_player INTEGER,
    winner VARCHAR(50),
    board_size INTEGER,
    number_of_boards INTEGER,
    difficulty INTEGER,
    game_history TEXT[][][], -- 3D array of characters
    gameover BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (session_id) REFERENCES Session(session_id) ON DELETE CASCADE
);

-- Create Wallet table
CREATE TABLE Wallet (
    uid INTEGER PRIMARY KEY,
    coins INTEGER DEFAULT 0,
    xp INTEGER DEFAULT 0,
    FOREIGN KEY (uid) REFERENCES Player(uid) ON DELETE CASCADE
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS Wallet;
DROP TABLE IF EXISTS SessionState;
DROP TABLE IF EXISTS Session;
DROP TABLE IF EXISTS "Player";
-- +goose StatementEnd


