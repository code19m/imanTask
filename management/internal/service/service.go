package service

import (
	"database/sql"
)

type service struct {
	db *sql.DB
}

func New(db *sql.DB) *service {
	return &service{
		db: db,
	}
}
