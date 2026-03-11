package main

import (
	"context"
	"fmt"
	"github.com/Royal17x/hireradar/internal/client/hh"
	"github.com/Royal17x/hireradar/internal/config"
	"github.com/Royal17x/hireradar/internal/delivery/bot"
	server "github.com/Royal17x/hireradar/internal/delivery/http"
	pg "github.com/Royal17x/hireradar/internal/repository/postgres"
	rd "github.com/Royal17x/hireradar/internal/repository/redis"
	"github.com/Royal17x/hireradar/internal/scheduler"
	"github.com/Royal17x/hireradar/internal/usecase"
	"github.com/Royal17x/hireradar/internal/utils"
	logger "github.com/charmbracelet/log"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Init logger
	logger.SetLevel(logger.DebugLevel)
	logger.SetReportCaller(true)

	// Load env variables
	if err := godotenv.Load(); err != nil {
		logger.Warn("Not found .env file")
	}

	// Load Conf
	cfg := config.MustLoad()

	// Init ctx
	ctx, cancel := context.WithCancel(context.Background())

	// Postgres Conn
	dbPool, err := pgxpool.New(ctx, cfg.Postgres.DSN())
	if err != nil {
		panic(err)
	}
	if err = dbPool.Ping(ctx); err != nil {
		panic(err)
	}
	logger.Info("Connected to PostgreSQL")
	defer dbPool.Close()

	// Migrations
	if err = utils.RunMigrations(cfg.Postgres.DSN()); err != nil {
		logger.Error("Migration error", "err", err)
	}

	// Redis Conn
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.Database,
	})
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	logger.Info("Connected to Redis", "response", pong)
	defer rdb.Close()

	// Repositories
	vacancyRepo := pg.NewVacancyRepository(dbPool)
	cacheRepo := rd.NewVacancyCache(rdb)
	userRepo := pg.NewUserRepository(dbPool)
	filterRepo := pg.NewFilterRepo(dbPool)
	accountRepo := pg.NewAccountRepo(dbPool)
	favoriteRepo := pg.NewFavoriteRepo(dbPool)

	// Client
	client := hh.New()

	// Usecase
	ucase := usecase.NewVacancyUsecase(vacancyRepo, cacheRepo, filterRepo, client)

	// Scheduler
	s := scheduler.NewScheduler(ucase, cfg.Parser.Interval, cfg.Parser.Query)
	go s.Start(ctx)

	// Bot
	tgBot, err := bot.NewBot(cfg.Telegram.Token, ucase, userRepo, filterRepo)
	if err != nil {
		logger.Error("Bot connection error", "err", err)
		os.Exit(1)
	}
	go tgBot.Start()

	//Server
	serv := server.NewServer(ucase, userRepo, filterRepo, accountRepo, favoriteRepo, cfg.Parser.Query, cfg.JWT.Secret)
	go func() {
		if err := serv.Run(":8080"); err != nil {
			logger.Error("Server run error", "err", err)
		}
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	logger.Info("Bot started")
	logger.Info("Server started")

	<-quit
	cancel()

	logger.Info("Gracefully shutting down...")
}
