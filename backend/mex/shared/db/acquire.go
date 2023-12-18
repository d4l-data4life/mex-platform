package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type keyType string

const ContextKeyTx keyType = "tx"

// Instead of just using `svc.DB` for accessing the database, we may also find a transaction
// in the context which we need to use (because the endpoint is called as part of a larger transaction).
// This function makes this distinction so that the caller can just normally use the transaction handle.
func AcquireTx(ctx context.Context, db *pgxpool.Pool) (pgx.Tx, error) {
	x := ctx.Value(ContextKeyTx)
	if x == nil {
		tx, err := db.BeginTx(ctx, pgx.TxOptions{
			IsoLevel: pgx.ReadCommitted, // Default, but let's make it explicit
		})
		if err != nil {
			return nil, err
		}
		return tx, nil
	}

	newTx, err := x.(pgx.Tx).Begin(ctx)
	if err != nil {
		return nil, err
	}
	return newTx, nil
}
