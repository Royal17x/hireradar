-- +goose Up
CREATE TABLE users(
    user_id SERIAL PRIMARY KEY,
    telegram_id TEXT UNIQUE NOT NULL,
    username VARCHAR(30) UNIQUE,
    is_active BOOLEAN NOT NULL ,
    created_at TIMESTAMP DEFAULT NOW()
);
-- +goose Down
DROP TABLE users;