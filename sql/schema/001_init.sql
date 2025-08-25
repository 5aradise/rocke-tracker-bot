-- +goose up
-- +goose statementbegin
CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    telegram_id INTEGER UNIQUE,
    email TEXT UNIQUE,
    subscription_destination TEXT NOT NULL CHECK (
        subscription_destination IN ('telegram', 'email')
    )
);

CREATE TABLE subscriptions (
    id INTEGER PRIMARY KEY,
    user_id INTEGER NOT NULL,
    players TEXT NOT NULL, -- '2x2', '3x3'
    mode TEXT NOT NULL, -- 'soccer', 'pentathlon'
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(user_id, players, mode)
);
-- +goose statementend

-- +goose down
-- +goose statementbegin
DROP TABLE IF EXISTS users;

DROP TABLE IF EXISTS subscriptions;
-- +goose statementend
