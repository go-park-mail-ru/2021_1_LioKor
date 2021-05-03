package common

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type PostgresDataBase struct {
	DBConn *pgxpool.Pool
}

func NewPostgresDataBase(cfg Config) (PostgresDataBase, error) {
	config, _ := pgxpool.ParseConfig("sslmode=disable")
	config.ConnConfig.Host = cfg.DBHost
	config.ConnConfig.Port = cfg.DBPort
	config.ConnConfig.User = cfg.DBUser
	config.ConnConfig.Password = cfg.DBPassword
	config.ConnConfig.Database = cfg.DBDatabase
	config.ConnConfig.ConnectTimeout = time.Duration(cfg.DBConnectTimeout) * time.Second
	dbpool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		return PostgresDataBase{}, err
	}

	return PostgresDataBase{dbpool}, nil
}

func (db *PostgresDataBase) Close() {
	db.DBConn.Close()
}
