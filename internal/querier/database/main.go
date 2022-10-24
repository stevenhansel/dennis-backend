package database

import "github.com/jmoiron/sqlx"

type DatabaseQuerier struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *DatabaseQuerier {
	return &DatabaseQuerier{
		db: db,
	}
}
