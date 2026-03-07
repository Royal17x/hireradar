package benchmarks

import (
	"context"
	"github.com/Royal17x/hireradar/internal/domain"
	"github.com/Royal17x/hireradar/internal/repository/postgres"
	redisRepo "github.com/Royal17x/hireradar/internal/repository/redis"
	"github.com/Royal17x/hireradar/internal/utils"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/redis/go-redis/v9"
	"github.com/testcontainers/testcontainers-go"
	pg "github.com/testcontainers/testcontainers-go/modules/postgres"
	rd "github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"os"
	"strings"
	"testing"
	"time"
)

var (
	vacancyRepo *postgres.VacancyRepository
	cacheRepo   *redisRepo.VacancyCache
)

func TestMain(m *testing.M) {
	os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")
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
	if err != nil {
		log.Fatal(err)
	}

	dsn, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	err = utils.RunMigrations(dsn)
	if err != nil {
		log.Fatal(err)
	}

	dbPool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatal(err)
	}

	vacancyRepo = postgres.NewVacancyRepository(dbPool)
	testVacancy := domain.Vacancy{
		HhID:        "test-hh-id",
		Title:       "test vacancy",
		Company:     "test company",
		URL:         "https://www.example.com",
		PublishedAt: time.Now(),
		CreatedAt:   time.Now(),
	}
	err = vacancyRepo.Save(ctx, &testVacancy)
	if err != nil {
		log.Fatal(err)
	}

	rdContainer, err := rd.Run(ctx,
		"redis:7-alpine",
		testcontainers.WithWaitStrategy(
			wait.ForListeningPort("6379/tcp")),
	)
	if err != nil {
		log.Fatal(err)
	}

	dsn, err = rdContainer.ConnectionString(ctx)
	if err != nil {
		log.Fatal(err)
	}

	dsn = strings.TrimPrefix(dsn, "redis://")

	redisClient := redis.NewClient(&redis.Options{
		Addr: dsn,
	})
	_, err = redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatal(err)
	}

	cacheRepo = redisRepo.NewVacancyCache(redisClient)

	err = cacheRepo.SetSeen(ctx, "test-hh-id")
	if err != nil {
		log.Fatal(err)
	}

	code := m.Run()
	pgContainer.Terminate(ctx)
	rdContainer.Terminate(ctx)
	redisClient.Close()
	dbPool.Close()
	os.Exit(code)
}
