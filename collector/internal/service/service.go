package service

import (
	"database/sql"

	"github.com/hibiken/asynq"
)

type service struct {
	taskClient *asynq.Client
	inspector  *asynq.Inspector
	db         *sql.DB
}

func New(redisClient *asynq.Client, inspector *asynq.Inspector, db *sql.DB) *service {
	return &service{
		taskClient: redisClient,
		inspector:  inspector,
		db:         db,
	}
}
