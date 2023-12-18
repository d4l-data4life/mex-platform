package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	L "github.com/d4l-data4life/mex/mex/shared/log"
)

type MigrationsProvider interface {
	InitScript() string
	ScriptForVersion(version int) string
}

func Migrate(ctx context.Context, log L.Logger, pool *pgxpool.Pool, migrationProvider MigrationsProvider) error {
	// Acquire connection and continue to use it throughout the function.
	// The advisory lock we use below only works when locked/unlocked from the same connection!
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	// Random but constant ID so that concurrent requests (e.g. from two starting replicas)
	// would try to acquire the lock for the same ID and just one of them prevails.
	lockSessionID := 0xa110ca7e

	start := time.Now()
	log.Info(ctx, L.Messagef("acquiring advisory lock (%d)", lockSessionID), L.Phase("migrations"))
	_, err = conn.Exec(ctx, "SELECT pg_advisory_lock($1)", lockSessionID)
	log.Info(ctx, L.Messagef("advisory lock acquisition duration: %v", time.Since(start)), L.Phase("migrations"))
	if err != nil {
		return err
	}
	defer func() {
		log.Info(ctx, L.Messagef("releasing advisory lock (%d)", lockSessionID), L.Phase("migrations"))
		_, err = conn.Exec(ctx, "SELECT pg_advisory_unlock($1)", lockSessionID)
		if err != nil {
			log.Warn(ctx, L.Messagef("could not unlock advisory lock %d (%s)", lockSessionID, err.Error()), L.Phase("migrations"))
		}
		log.Info(ctx, L.Messagef("advisory lock released (%d)", lockSessionID), L.Phase("migrations"))
	}()

	// Arriving here means we have an advisory lock and can do the migration(s).

	nextMigrationScript := ""
	i := 0
	for {
		i++
		if i > 64 {
			return fmt.Errorf("too many migrations (%d)", i)
		}

		nextVersion, err := getNextMigrationVersion(ctx, conn)
		if err != nil {
			nextMigrationScript = migrationProvider.InitScript()
		} else {
			nextMigrationScript = migrationProvider.ScriptForVersion(nextVersion)
			if nextMigrationScript == "" {
				log.Info(ctx, L.Messagef("no further migration script: %d; migration done/unnecessary", nextVersion), L.Phase("migrations"))
				return nil
			}
		}

		log.Info(ctx, L.Messagef("run migration: %d", nextVersion), L.Phase("migrations"))
		err = runMigration(ctx, conn, nextMigrationScript)
		if err != nil {
			log.Error(ctx, L.Messagef(err.Error()))
			return nil
		}
	}
}

func getNextMigrationVersion(ctx context.Context, conn *pgxpool.Conn) (int, error) {
	rows, err := conn.Query(ctx, "SELECT * FROM next_migration_version()")
	if err != nil {
		return -1, err
	}
	defer rows.Close()

	rows.Next()
	var nextMigrationVersion int
	err = rows.Scan(&nextMigrationVersion)
	if err != nil {
		return -1, err
	}

	return nextMigrationVersion, nil
}

func runMigration(ctx context.Context, conn *pgxpool.Conn, script string) error {
	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}

	_, err = conn.Exec(ctx, script)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	_ = tx.Commit(ctx)
	return nil
}
