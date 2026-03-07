package postgres

import (
	"context"
	"github.com/Royal17x/hireradar/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FilterRepo struct {
	db *pgxpool.Pool
}

func NewFilterRepo(db *pgxpool.Pool) *FilterRepo {
	return &FilterRepo{db: db}
}

func (f *FilterRepo) Save(ctx context.Context, filter domain.Filter) error {
	query := `
INSERT INTO filters (user_id, keywords, city, grade)
VALUES ($1, $2, $3, $4)`
	_, err := f.db.Exec(ctx, query, filter.UserID, filter.Keywords, filter.City, filter.Grade)
	if err != nil {
		return err
	}
	return nil
}

func (f *FilterRepo) GetByUserID(ctx context.Context, userID int) ([]domain.Filter, error) {
	query := `
SELECT id, user_id, keywords, city, grade, created_at
FROM filters
WHERE user_id = $1`
	rows, err := f.db.Query(ctx, query, &userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var filters []domain.Filter
	for rows.Next() {
		var filter domain.Filter
		rows.Scan(&filter.ID, &filter.UserID, &filter.Keywords, &filter.City, &filter.Grade, &filter.CreatedAt)
		filters = append(filters, filter)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return filters, nil
}

func (f *FilterRepo) Update(ctx context.Context, filter domain.Filter) error {
	query := `UPDATE filters
SET keywords = $1, city = $2, grade = $3
WHERE id = $4`
	_, err := f.db.Exec(ctx, query, filter.Keywords, filter.City, filter.Grade, filter.ID)
	if err != nil {
		return err
	}
	return nil
}

func (f *FilterRepo) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM filters WHERE id = $1`
	_, err := f.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}
