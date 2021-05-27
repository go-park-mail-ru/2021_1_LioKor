package common

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"fmt"
)

type GormPostgresDataBase struct{
	DB *gorm.DB
}

func NewGormPostgresDataBase(cfg Config)(GormPostgresDataBase, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBDatabase,
		cfg.DBPort,
		)
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: dsn,
		PreferSimpleProtocol: true,
	},
	), &gorm.Config{})
	if err != nil {
		return GormPostgresDataBase{}, err
	}
	postgresDB, err := db.DB()
	if err != nil {
		return GormPostgresDataBase{}, err
	}
	err = postgresDB.Ping()
	if err != nil {
		postgresDB.Close()
		return GormPostgresDataBase{}, err
	}
	postgresDB.SetMaxOpenConns(50)

	return GormPostgresDataBase{
		DB: db,
	}, nil
}

func (db *GormPostgresDataBase) Close() {
	sqlDb, _ := db.DB.DB()
	sqlDb.Close()
}

func (db *GormPostgresDataBase) AddMail() {

}
