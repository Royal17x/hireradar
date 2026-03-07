package postgres

import (
	"context"
	"errors"
	"github.com/Royal17x/hireradar/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Save(ctx context.Context, user domain.User) error {
	query := `
INSERT INTO users (telegram_id, username, is_active, created_at)
VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(ctx, query, user.TgID, user.Username, user.IsActive, user.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) GetByTelegramID(ctx context.Context, tgID string) (*domain.User, error) {
	query := `SELECT user_id, telegram_id, username, is_active, created_at
FROM users 
WHERE telegram_id=$1;`
	user := domain.User{}
	err := r.db.QueryRow(ctx, query, tgID).Scan(&user.UserID, &user.TgID, &user.Username, &user.IsActive, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
