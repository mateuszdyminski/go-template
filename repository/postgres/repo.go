package postgres

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/mateuszdyminski/go-template/app"
)

type pgRepository struct {
	db *sql.DB
}

// NewPostgresRepository - returns new repository which connects to postgres DB.
// It implements app.Repository interface.
func NewPostgresRepository(host string, port int, user, password, dbname string) (app.Repository, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, errors.Wrap(err, "can't create postgres repo")
	}

	return &pgRepository{db: db}, nil
}

// OK - returns information whether connection to DB is up and running.
func (r *pgRepository) OK(ctx context.Context) (bool, error) {
	if err := r.db.PingContext(ctx); err != nil {
		return false, errors.Wrap(err, "postgres ping failed")
	}

	return true, nil
}
