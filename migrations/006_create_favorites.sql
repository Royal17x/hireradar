-- +goose Up
CREATE TABLE favorites(
    id SERIAL PRIMARY KEY,
    account_id INTEGER NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    hh_id TEXT NOT NULL,
    UNIQUE (account_id, hh_id)
);
-- +goose Down
DROP TABLE favorites;
