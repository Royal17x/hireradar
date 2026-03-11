package postgres

import (
	"context"
	"errors"
	"github.com/Royal17x/hireradar/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AccountRepo struct {
	dbPool *pgxpool.Pool
}

func NewAccountRepo(dbPool *pgxpool.Pool) *AccountRepo {
	return &AccountRepo{dbPool: dbPool}
}

func (a *AccountRepo) Save(ctx context.Context, account domain.Account) (int, error) {
	query := `INSERT INTO accounts (email, password_hash, created_at)
VALUES ($1,$2,$3) RETURNING id;`
	var id int
	err := a.dbPool.QueryRow(ctx, query, account.Email, account.PasswordHash, account.CreatedAt).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, err
}

func (a *AccountRepo) GetByEmail(ctx context.Context, email string) (*domain.Account, error) {
	query := `SELECT id, email, password_hash, created_at 
FROM accounts
WHERE email=$1;`
	account := &domain.Account{}
	err := a.dbPool.QueryRow(ctx, query, email).Scan(&account.ID, &account.Email, &account.PasswordHash, &account.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return account, nil
}

func (a *AccountRepo) GetByID(ctx context.Context, id int) (*domain.Account, error) {
	query := `SELECT id, email, password_hash, created_at
FROM accounts
WHERE id=$1;`
	account := &domain.Account{}
	err := a.dbPool.QueryRow(ctx, query, id).Scan(&account.ID, &account.Email, &account.PasswordHash, &account.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return account, nil
}
