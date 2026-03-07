-- +goose Up
ALTER TABLE vacancies ADD COLUMN city VARCHAR(100);

-- +goose Down
ALTER TABLE vacancies DROP COLUMN city;