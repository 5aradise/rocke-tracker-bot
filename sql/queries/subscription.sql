-- name: CreateSubscriptionByTelegramID :execlastid
INSERT INTO subscriptions (user_id, players, mode)
VALUES (
    (SELECT id FROM users WHERE telegram_id = ?), ?, ?
);

-- name: ListSubscriptionsByTelegramID :many
SELECT sub.id, sub.players, sub.mode FROM subscriptions sub
JOIN users ON users.id = sub.user_id
WHERE users.telegram_id = ?;

-- name: ListTelegramIDsBySubscription :many
SELECT telegram_id FROM users
JOIN subscriptions ON subscriptions.user_id = users.id
WHERE subscriptions.players = ?
AND subscriptions.mode = ?;

-- name: DeleteSubscription :execrows
DELETE FROM subscriptions
WHERE user_id = ( SELECT id FROM users WHERE telegram_id = ? )
AND players = ?
AND mode = ?;
