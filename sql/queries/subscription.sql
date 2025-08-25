-- name: CreateSubscriptionByTelegramID :execlastid
INSERT INTO subscriptions (user_id, players, mode)
VALUES (
    (SELECT id FROM users WHERE telegram_id = ?), ?, ?
);

-- name: DeleteSubscription :execrows
DELETE FROM subscriptions
WHERE user_id = ( SELECT id FROM users WHERE telegram_id = ? )
AND players = ?
AND mode = ?;
