package common

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"context"
)


type PostgresDataBase struct {
	DBConn *pgxpool.Conn
	DBpool *pgxpool.Pool
}

func NewPostgresDataBase(dbConfig string) (PostgresDataBase, error){
	dbpool, err := pgxpool.Connect(context.Background(), dbConfig)
	if err != nil {
		return PostgresDataBase{}, err
	}

	conn, err := dbpool.Acquire(context.Background())
	if err != nil {
		return PostgresDataBase{}, err
	}

	return PostgresDataBase{conn, dbpool}, nil
}

func (db *PostgresDataBase) Close() {
	db.DBConn.Release()
	db.DBpool.Close()
}