-- name: CreateUser :execlastid
INSERT INTO users (telegram_id, subscription_destination) 
VALUES (?, 'telegram');

-- name: DeleteUserByTelegramID :execrows
DELETE FROM users
WHERE telegram_id = ?;
