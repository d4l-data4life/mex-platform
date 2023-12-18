package pglo

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	L "github.com/d4l-data4life/mex/mex/shared/log"
)

// Postgres Large Object blob store (https://www.postgresql.org/docs/current/largeobjects.html)
type PostgresLargeObjectStore struct {
	DB              *pgxpool.Pool
	MasterTableName string
	Log             L.Logger
}

func (store *PostgresLargeObjectStore) logInfo(ctx context.Context, opts ...L.Opt) {
	if store.Log == nil {
		return
	}
	store.Log.Info(ctx, opts...)
}

func finishTx(ctx context.Context, tx pgx.Tx, log L.Logger, commit *bool) {
	if *commit {
		if err := tx.Commit(ctx); err != nil {
			if log != nil {
				log.Error(ctx, L.Messagef("commit error: %s", err.Error()))
			}
		}
	} else {
		if err := tx.Rollback(ctx); err != nil {
			if log != nil {
				log.Error(ctx, L.Messagef("rollback error: %s", err.Error()))
			}
		}
	}
}
