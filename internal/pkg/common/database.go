package common

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PostgresDataBase struct {
	DBConn *pgxpool.Pool
}

func NewPostgresDataBase(dbConfig string) (PostgresDataBase, error) {
	dbpool, err := pgxpool.Connect(context.Background(), dbConfig)
	if err != nil {
		return PostgresDataBase{}, err
	}

	return PostgresDataBase{dbpool}, nil
}

func (db *PostgresDataBase) Close() {
	db.DBConn.Close()
}
