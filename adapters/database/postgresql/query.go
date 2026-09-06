// Package postgresql implements the core tank, pump, and volume repository
// interfaces against a PostgreSQL database, using pgx for connections and
// scany to scan rows into domain structs.
package postgresql

import (
	"context"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

// archimedesPool wraps a pgxpool.Pool with query helpers shared by every
// repository in this package.
type archimedesPool struct {
	*pgxpool.Pool
}

// NewPool wraps an existing pgxpool.Pool for use by the repositories in
// this package. The pool's lifecycle (creation, closing) is owned by the
// caller.
func NewPool(pgxpool *pgxpool.Pool) *archimedesPool {
	return &archimedesPool{
		Pool: pgxpool,
	}
}

// ReadArchimedes acquires a connection, runs query, and scans every
// resulting row into dst (a pointer to a slice or struct, per scany's
// conventions).
func (p *archimedesPool) ReadArchimedes(ctx context.Context, dst any, query string, args ...any) error {
	connection, err := p.Acquire(ctx)
	if err != nil {
		return err
	}
	defer connection.Release()

	rows, err := connection.Query(ctx, query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	err = pgxscan.ScanAll(dst, rows)
	if err != nil {
		return err
	}
	return nil
}

// WriteArchimedes acquires a connection and runs query for its side
// effects, reporting anything the database rejected it for: a constraint
// violation, a missing relation, a trigger raising.
//
// It uses Exec rather than Query even for the RETURNING clauses the write
// statements carry. Query reports only the errors raised before execution
// begins, and surfaces the rest through the rows — so a write whose rows
// are never read comes back as a success, and an event the database
// refused would be recorded nowhere and logged by no one.
//
// What it still does not report is a write that was accepted and matched
// nothing: an UPDATE/INSERT whose WHERE clause selects no row succeeds
// silently rather than returning a not-found error, which is what lets
// StopPump ignore a pump with no open run.
func (p *archimedesPool) WriteArchimedes(ctx context.Context, query string, args ...any) error {
	connection, err := p.Acquire(ctx)
	if err != nil {
		return err
	}
	defer connection.Release()

	_, err = connection.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
