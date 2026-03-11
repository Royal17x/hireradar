package postgres

import (
	"context"
	"github.com/Royal17x/hireradar/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FavoriteRepo struct {
	dbPool *pgxpool.Pool
}

func NewFavoriteRepo(dbPool *pgxpool.Pool) *FavoriteRepo {
	return &FavoriteRepo{dbPool: dbPool}
}

func (f *FavoriteRepo) Save(ctx context.Context, favorite domain.Favorite) error {
	query := `INSERT INTO favorites (account_id, hh_id) VALUES ($1, $2)`
	_, err := f.dbPool.Exec(ctx, query, favorite.AccountID, favorite.HhID)
	if err != nil {
		return err
	}
	return nil
}

func (f *FavoriteRepo) GetByAccountID(ctx context.Context, accountID int) ([]domain.Favorite, error) {
	query := `SELECT id, account_id, hh_id FROM favorites WHERE account_id = $1`
	rows, err := f.dbPool.Query(ctx, query, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var favorites []domain.Favorite
	for rows.Next() {
		var favorite domain.Favorite
		err := rows.Scan(&favorite.ID, &favorite.AccountID, &favorite.HhID)
		if err != nil {
			return nil, err
		}
		favorites = append(favorites, favorite)
	}
	if rows.Err() != nil {
		return nil, err
	}
	return favorites, nil
}

func (f *FavoriteRepo) Delete(ctx context.Context, accountID int, hhID string) error {
	query := `DELETE FROM favorites WHERE account_id = $1 AND hh_id = $2`
	_, err := f.dbPool.Exec(ctx, query, accountID, hhID)
	if err != nil {
		return err
	}
	return nil
}
