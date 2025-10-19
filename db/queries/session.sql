-- name: GetLatestSessionStateByPlayerId :one
SELECT 
    s.session_id,
    s.uid,
    s.expired,
    s.created_at,
    ss.boards,
    ss.current_player,
    ss.winner,
    ss.board_size,
    ss.number_of_boards,
    ss.difficulty,
    ss.game_history,
    ss.gameover
FROM session s
JOIN sessionstate ss
  ON s.session_id = ss.session_id
WHERE s.uid = $1
  AND s.created_at >= now() - interval '15 minutes'
ORDER BY s.created_at DESC
LIMIT 1;

-- name: CreateSession :exec
INSERT INTO session (session_id, uid, expired, created_at)
VALUES ($1, $2, false, now());
-- name: CreateInitialSessionState :exec
INSERT INTO sessionstate (
    session_id,
    boards,
    current_player,
    winner,
    board_size,
    number_of_boards,
    difficulty,
    game_history,
    gameover
)
VALUES (
    $1,                          
    $2,                           
    1,                            
    '',                           
    $3,                           
    $4,                           
    $5,                           
    $6,                  
    false                         
);