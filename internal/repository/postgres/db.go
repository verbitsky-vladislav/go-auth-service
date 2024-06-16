package postgres

import (
	"auth-microservice/pkg/logger"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewDB(host, port, user, password, dbname string) (*sqlx.DB, error) {
	dbURL := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		return nil, logger.Error(err, "failed to connect to postgres")
	}

	if err = db.Ping(); err != nil {
		return nil, logger.Error(err, "failed to ping postgres")
	}

	return db, nil
}
