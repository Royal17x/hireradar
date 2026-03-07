package main

import (
	"context"
	"fmt"
	"github.com/Royal17x/hireradar/internal/client/hh"
	"github.com/Royal17x/hireradar/internal/config"
	"github.com/Royal17x/hireradar/internal/delivery/bot"
	pg "github.com/Royal17x/hireradar/internal/repository/postgres"
	rd "github.com/Royal17x/hireradar/internal/repository/redis"
	"github.com/Royal17x/hireradar/internal/scheduler"
	"github.com/Royal17x/hireradar/internal/usecase"
	"github.com/Royal17x/hireradar/internal/utils"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	//Load env variables
	if err := godotenv.Load(); err != nil {
		log.Println("не найден .env файл с env variables")
	}

	//Load Conf
	cfg := config.MustLoad()

	//Init ctx
	ctx, cancel := context.WithCancel(context.Background())

	//Postgres Conn
	dbPool, err := pgxpool.New(ctx, cfg.Postgres.DSN())
	if err != nil {
		panic(err)
	}
	if err = dbPool.Ping(ctx); err != nil {
		panic(err)
	}
	fmt.Println("подключились к postgres")
	defer dbPool.Close()

	// migrations
	if err = utils.RunMigrations(cfg.Postgres.DSN()); err != nil {
		log.Fatalf("ошибка миграций: %v", err)
	}

	//Redis Conn
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.Database,
	})
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("подклюлись к redis", pong)
	defer rdb.Close()

	//repositories
	vacancyRepo := pg.NewVacancyRepository(dbPool)
	cacheRepo := rd.NewVacancyCache(rdb)
	userRepo := pg.NewUserRepository(dbPool)
	filterRepo := pg.NewFilterRepo(dbPool)

	//client
	client := hh.New()

	//usecase
	ucase := usecase.NewVacancyUsecase(vacancyRepo, cacheRepo, client)

	//scheduler
	s := scheduler.NewScheduler(ucase, cfg.Parser.Interval, "golang")
	go s.Start(ctx)

	//bot
	tgBot, err := bot.NewBot(cfg.Telegram.Token, ucase, userRepo, filterRepo)
	if err != nil {
		log.Fatal(err)
	}
	go tgBot.Start()

	//Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	log.Println("бот запущен")

	<-quit
	cancel()

	log.Println("gracefully завершаем")
}
