package main

import (
	"context"
	"database/sql"
	"management/internal/config"
	"management/internal/handler/grpcsrv"
	"management/internal/service"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	service := service.New(db)

	grpcSrv := grpcsrv.New(service)

	go grpcSrv.Run(cfg.GrpcServerAddr)

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = 10 * time.Second

	_, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	grpcSrv.ShutDown()

	err = db.Close()
	if err != nil {
		log.Fatal().Err(err).Msg("cannot gracefully close db conn")
	}

	log.Info().Msg("Application gracefully shut down...")
}
