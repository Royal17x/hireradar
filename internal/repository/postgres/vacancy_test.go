package postgres

import (
	"context"
	"github.com/Royal17x/hireradar/internal/domain"
	"github.com/Royal17x/hireradar/internal/utils"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	pg "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"testing"
	"time"
)

func TestVacancyRepository(t *testing.T) {
	t.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")
	ctx := context.Background()
	pgContainer, err := pg.Run(ctx,
		"postgres:16-alpine",
		pg.WithDatabase("testdb"),
		pg.WithUsername("testuser"),
		pg.WithPassword("testpassword"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2)),
	)
	require.NoError(t, err)
	defer pgContainer.Terminate(ctx)

	dsn, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	err = utils.RunMigrations(dsn)
	require.NoError(t, err)

	dbPool, err := pgxpool.New(ctx, dsn)
	require.NoError(t, err)
	defer dbPool.Close()

	vacancyRepo := NewVacancyRepository(dbPool)

	t.Run("save - сохраняет вакансию", func(t *testing.T) {
		vacancy := domain.Vacancy{
			HhID:        "123",
			Title:       "golang dev",
			Company:     "yandex",
			URL:         "https://golang.org",
			PublishedAt: time.Now(),
			CreatedAt:   time.Now(),
		}

		err := vacancyRepo.Save(ctx, &vacancy)
		require.NoError(t, err)

	})

	t.Run("save - дубль не сохраняется", func(t *testing.T) {
		vacancy := domain.Vacancy{
			HhID:        "123",
			Title:       "golang dev",
			Company:     "yandex",
			URL:         "https://golang.org",
			PublishedAt: time.Now(),
			CreatedAt:   time.Now(),
		}
		err := vacancyRepo.Save(ctx, &vacancy)
		require.Error(t, err)

	})

	t.Run("getAll - возвращает сохранённые вакансии", func(t *testing.T) {
		vacancies, err := vacancyRepo.GetAll(ctx)
		require.NoError(t, err)
		require.True(t, len(vacancies) > 0)
		found := false
		for _, v := range vacancies {
			if v.HhID == "123" {
				found = true
			}
		}
		require.True(t, found)
	})

	t.Run("exists - возвращает true для существующей", func(t *testing.T) {
		exists, err := vacancyRepo.Exists(ctx, "123")
		require.NoError(t, err)
		require.True(t, exists)
	})
}
