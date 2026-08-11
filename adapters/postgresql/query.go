package postgresql

import (
	"context"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

type archimedesPool struct {
	*pgxpool.Pool
}

func NewPool(pgxpool *pgxpool.Pool) *archimedesPool {
	return &archimedesPool{
		Pool: pgxpool,
	}
}

func (p *archimedesPool) QueryArchimedes(ctx context.Context, dst any, query string, args ...any) error {
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
