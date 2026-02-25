package main

import (
	"context"
	"fmt"
	"github.com/Royal17x/hireradar/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
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

	//TODO: parser here + start bot

	//Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	log.Println("бот запущен")

	<-quit
	cancel()

	log.Println("gracefully завершаем")
}
