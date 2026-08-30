-- name: UpsertImageForUID :one
INSERT INTO images (file_id, uid, file_path)
VALUES ($1, $2, $3)
ON CONFLICT (file_id) DO UPDATE
SET file_path = EXCLUDED.file_path,
    updated_at = NOW()
WHERE images.uid = EXCLUDED.uid
RETURNING file_id, uid, file_path, source, created_at, updated_at;
