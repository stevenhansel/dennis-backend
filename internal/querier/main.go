package querier

import "github.com/jmoiron/sqlx"

type Querier struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Querier {
	return &Querier{
		db: db,
	}
}
