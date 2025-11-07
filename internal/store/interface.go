package store

import (
	"context"
)

type Database interface {
	Exec(ctx context.Context, query string, args ...interface{}) error
	Query(ctx context.Context, query string, args ...interface{}) (Rows, error)
	QueryRow(ctx context.Context, query string, args ...interface{}) Row
	BeginTx(ctx context.Context) (Transaction, error)
	Close() error
	Ping(ctx context.Context) error
	RunMigrations() error

	Connect(ctx context.Context) error
}

type Transaction interface {
	Exec(ctx context.Context, query string, args ...interface{}) error
	Query(ctx context.Context, query string, args ...interface{}) (Rows, error)
	QueryRow(ctx context.Context, query string, args ...interface{}) Row
	Commit() error
	Rollback() error
}

type Rows interface {
	Close() error
	Next() bool
	Scan(dest ...interface{}) error
	Err() error
}

type Row interface {
	Scan(dest ...interface{}) error
}