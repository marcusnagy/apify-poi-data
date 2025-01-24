package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	Pool    *pgxpool.Pool
	Queries *Queries
}

func NewDatabase(ctx context.Context, connString string) (*Database, error) {
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	queries := New(pool)

	return &Database{
		Pool:    pool,
		Queries: queries,
	}, nil
}

func (d *Database) Close() {
	d.Pool.Close()
}
