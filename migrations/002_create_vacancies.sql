-- +goose Up
CREATE TABLE vacancies(
  vacancy_id SERIAL PRIMARY KEY,
  hh_id TEXT UNIQUE NOT NULL,
  title TEXT NOT NULL,
  company TEXT NOT NULL,
  url TEXT NOT NULL,
  salary_from INTEGER,
  salary_to INTEGER,
  published_at TIMESTAMP NOT NULL,
  created_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX idx_vacancy_title ON vacancies
USING GIN (to_tsvector('russian', title));
-- +goose Down