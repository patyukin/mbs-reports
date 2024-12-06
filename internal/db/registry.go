package db

import (
	"context"
	"database/sql"
)

type QueryExecutor interface {
	ExecContext(ctx context.Context, q string, args ...interface{}) (sql.Result, error)
	Query(q string, args ...interface{}) (*sql.Rows, error)
}

type Registry struct {
	db *sql.DB
}

func New(db *sql.DB) *Registry {
	return &Registry{db: db}
}

func (registry *Registry) GetRepo() *Repository {
	return &Repository{
		db: registry.db,
	}
}

type Handler func(ctx context.Context, repo *Repository) error

func (registry *Registry) Close() error {
	return registry.db.Close()
}
