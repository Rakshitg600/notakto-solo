-- +goose Up
ALTER TABLE SessionState DROP COLUMN current_player;

-- +goose Down
ALTER TABLE SessionState ADD COLUMN current_player INTEGER;