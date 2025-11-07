package postgresql

import (
	"context"
	"database/sql"
	"denet/internal/store"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type PostgresDatabase struct {
	db    *sql.DB
	dbURL string
}

type PostgresTransaction struct {
	tx *sql.Tx
}

type PostgresRows struct {
	rows *sql.Rows
}

type PostgresRow struct {
	row *sql.Row
}

func NewPostgresDatabase(dbURL string) store.Database {
	return &PostgresDatabase{
		dbURL: dbURL,
	}
}

func (p *PostgresDatabase) Connect(ctx context.Context) error {
	db, err := sql.Open("postgres", p.dbURL)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	p.db = db
	return nil
}

func (p *PostgresDatabase) Close() error {
	if p.db != nil {
		return p.db.Close()
	}
	return nil
}

func (p *PostgresDatabase) Ping(ctx context.Context) error {
	if p.db == nil {
		return fmt.Errorf("database not connected")
	}
	return p.db.PingContext(ctx)
}

func (p *PostgresDatabase) Exec(ctx context.Context, query string, args ...interface{}) error {
	if p.db == nil {
		return fmt.Errorf("database not connected")
	}
	_, err := p.db.ExecContext(ctx, query, args...)
	return err
}

func (p *PostgresDatabase) Query(ctx context.Context, query string, args ...interface{}) (store.Rows, error) {
	if p.db == nil {
		return nil, fmt.Errorf("database not connected")
	}
	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &PostgresRows{rows: rows}, nil
}

func (p *PostgresDatabase) QueryRow(ctx context.Context, query string, args ...interface{}) store.Row {
	if p.db == nil {
		return &PostgresRow{}
	}
	row := p.db.QueryRowContext(ctx, query, args...)
	return &PostgresRow{row: row}
}

func (p *PostgresDatabase) BeginTx(ctx context.Context) (store.Transaction, error) {
	if p.db == nil {
		return nil, fmt.Errorf("database not connected")
	}
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &PostgresTransaction{tx: tx}, nil
}

func (p *PostgresDatabase) RunMigrations() error {
	if p.db == nil {
		return fmt.Errorf("database not connected")
	}
	driver, err := postgres.WithInstance(p.db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	log.Println("Migrations applied successfully")
	return nil
}

func (pt *PostgresTransaction) Exec(ctx context.Context, query string, args ...interface{}) error {
	_, err := pt.tx.ExecContext(ctx, query, args...)
	return err
}

func (pt *PostgresTransaction) Query(ctx context.Context, query string, args ...interface{}) (store.Rows, error) {
	rows, err := pt.tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &PostgresRows{rows: rows}, nil
}

func (pt *PostgresTransaction) QueryRow(ctx context.Context, query string, args ...interface{}) store.Row {
	row := pt.tx.QueryRowContext(ctx, query, args...)
	return &PostgresRow{row: row}
}

func (pt *PostgresTransaction) Commit() error {
	return pt.tx.Commit()
}

func (pt *PostgresTransaction) Rollback() error {
	return pt.tx.Rollback()
}

func (pr *PostgresRows) Close() error {
	return pr.rows.Close()
}

func (pr *PostgresRows) Next() bool {
	return pr.rows.Next()
}

func (pr *PostgresRows) Scan(dest ...interface{}) error {
	return pr.rows.Scan(dest...)
}

func (pr *PostgresRows) Err() error {
	return pr.rows.Err()
}

func (pr *PostgresRow) Scan(dest ...interface{}) error {
	if pr.row == nil {
		return fmt.Errorf("row is nil")
	}
	return pr.row.Scan(dest...)
}