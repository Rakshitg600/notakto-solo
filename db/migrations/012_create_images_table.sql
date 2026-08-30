-- +goose Up
-- +goose StatementBegin
CREATE TABLE images (
    file_id TEXT PRIMARY KEY,
    uid VARCHAR(36) NOT NULL,
    file_path TEXT NOT NULL,
    source TEXT NOT NULL DEFAULT 'Imagekit',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (uid) REFERENCES Player(uid) ON DELETE CASCADE
);

CREATE INDEX idx_images_uid ON images(uid);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_images_uid;
DROP TABLE IF EXISTS images;
-- +goose StatementEnd
