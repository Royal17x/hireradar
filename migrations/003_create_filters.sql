-- +goose Up
CREATE TABLE filters(
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    keywords TEXT NOT NULL,
    city VARCHAR(30),
    grade TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

CREATE INDEX idx_user_id ON filters (user_id);
CREATE INDEX idx_keywords ON filters
USING GIN (to_tsvector('russian', keywords));
-- +goose Down
