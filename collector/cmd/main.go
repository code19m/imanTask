package main

import (
	"collector/internal/config"
	"collector/internal/handler/grpcsrv"
	"collector/internal/handler/worker"
	"collector/internal/service"
	"collector/pkg/logger"
	"context"
	"database/sql"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"

	_ "modernc.org/sqlite"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
	}

	db, err := sql.Open("sqlite", "./posts.db")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot open connection to sqlite")
	}
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	createTableSql := `CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY, user_id INTEGER, title VARCHAR, body TEXT);`

	_, err = db.Exec(createTableSql)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot execute createTableSql")
	}

	redisOpt := asynq.RedisClientOpt{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Pass,
		DB:       cfg.Redis.DB,
	}

	redisCl := asynq.NewClient(redisOpt)
	inspector := asynq.NewInspector(redisOpt)

	service := service.New(redisCl, inspector, db)

	workerSrv := asynq.NewServer(redisOpt, asynq.Config{
		Queues: map[string]int{
			"default": 5,
		},
		ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
			log.Error().Err(err).Str("type", task.Type()).
				Bytes("payload", task.Payload()).Msg("process task failed")
		}),
		Logger: logger.NewLogger(),
	})

	worker := worker.New(workerSrv, service)
	grpcSrv := grpcsrv.New(service)

	go worker.Run()
	go grpcSrv.Run(cfg.GrpcServerAddr)

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = 10 * time.Second

	_, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	grpcSrv.ShutDown()
	worker.ShutDown()

	err = db.Close()
	if err != nil {
		log.Fatal().Err(err).Msg("cannot gracefully close db conn")
	}

	redisCl.Close()
	inspector.Close()

	log.Info().Msg("Application gracefully shut down...")
}
