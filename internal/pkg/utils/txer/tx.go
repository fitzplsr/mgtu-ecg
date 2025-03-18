package txer

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

//type runFunc func(ctx context.Context, tx pgx.Tx)

func Tx(ctx context.Context, conn *pgxpool.Pool) (pgx.Tx, error) {
	return conn.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.Serializable,
	})
}
