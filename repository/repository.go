// This file contains the repository implementation layer.
package repository

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

const (
	driverPostgresDB = "postgres"
)

type Repository struct {
	db *sql.DB
}

func New(dbDsn string) *Repository {
	// open mere validating connection string
	pgdb, err := sql.Open(driverPostgresDB, dbDsn)
	if err != nil {
		log.Panicf("[NewRepository] invalid connection string, err:%+v\n", err)
	}

	// validate that connection indeed can be established
	err = pgdb.Ping()
	if err != nil {
		log.Panicf("[NewRepository] failed to established connection with database, err:%+v\n", err)
	}

	return &Repository{
		db: pgdb,
	}
}
