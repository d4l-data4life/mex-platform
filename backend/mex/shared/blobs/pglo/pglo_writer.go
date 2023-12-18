package pglo

import (
	"context"
	"fmt"
	"io"

	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc/codes"

	"github.com/d4l-data4life/mex/mex/shared/errstat"
	L "github.com/d4l-data4life/mex/mex/shared/log"
)

const (
	sqlSelectBlobInfo = "SELECT blob_oid FROM %s WHERE blob_name = $1 AND blob_type = $2"
	sqlDeleteBlobInfo = "DELETE FROM %s WHERE blob_name = $1 AND blob_type = $2"
	sqlUpsertBlobInfo = `INSERT INTO %s (blob_name, blob_type, blob_oid)
						 VALUES ($1, $2, $3) ON CONFLICT(blob_name, blob_type) DO UPDATE set blob_oid = EXCLUDED.blob_oid`
)

func (store *PostgresLargeObjectStore) Delete(ctx context.Context, blobName string, blobType string) error {
	var err error

	if store.DB == nil {
		return fmt.Errorf("store DB is nil")
	}

	tx, err := store.DB.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("error beginning transaction: %w", err)
	}

	txCommit := false
	defer finishTx(ctx, tx, store.Log, &txCommit)

	var oid uint32
	err = tx.QueryRow(context.Background(), fmt.Sprintf(sqlSelectBlobInfo, store.MasterTableName), blobName, blobType).Scan(&oid)
	if err == pgx.ErrNoRows {
		return errstat.MakeGRPCStatus(codes.NotFound, fmt.Sprintf("blob not found: %s/%s", blobName, blobType)).Err()
	} else if err != nil {
		return fmt.Errorf("query error: %s/%s: %w", blobName, blobType, err)
	}

	lo := tx.LargeObjects()
	err = lo.Unlink(ctx, oid)
	if err != nil {
		return fmt.Errorf("error deleting blob: %s/%s: %w", blobName, blobType, err)
	}

	_, err = tx.Exec(ctx, fmt.Sprintf(sqlDeleteBlobInfo, store.MasterTableName), blobName, blobType)
	if err != nil {
		return fmt.Errorf("error deleting blob master data: %s/%s: %w", blobName, blobType, err)
	}

	txCommit = true
	return nil
}

func (store *PostgresLargeObjectStore) GetWriteCloser(ctx context.Context, blobName string, blobType string, append bool) (io.WriteCloser, error) {
	var err error

	if store.DB == nil {
		return nil, fmt.Errorf("store DB is nil")
	}

	tx, err := store.DB.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("error beginning transaction")
	}

	lo := tx.LargeObjects()

	var oid uint32
	err = store.DB.QueryRow(ctx, fmt.Sprintf(sqlSelectBlobInfo, store.MasterTableName), blobName, blobType).Scan(&oid)
	if err == pgx.ErrNoRows {
		oid, err = lo.Create(ctx, 0)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	store.logInfo(ctx, L.Messagef("OID: %d", oid))

	lob, err := lo.Open(ctx, oid, pgx.LargeObjectModeRead|pgx.LargeObjectModeWrite)
	if err != nil {
		return nil, err
	}

	if append {
		_, err = lob.Seek(0, io.SeekEnd)
	} else {
		err = lob.Truncate(0)
	}
	if err != nil {
		return nil, err
	}

	return &lobWriteWrapper{
		lob:      lob,
		txCommit: true,
		finisher: func(txCommit bool) {
			defer finishTx(ctx, tx, store.Log, &txCommit)
			if _, err := tx.Exec(context.Background(), fmt.Sprintf(sqlUpsertBlobInfo, store.MasterTableName), blobName, blobType, oid); err != nil {
				txCommit = false
			}
		},
	}, nil
}

type lobWriteWrapper struct {
	lob      *pgx.LargeObject
	txCommit bool
	finisher func(bool)
}

func (lobw *lobWriteWrapper) Write(data []byte) (int, error) {
	if !lobw.txCommit {
		return -1, fmt.Errorf("previous write failed, ignoring future writes")
	}

	n, err := lobw.lob.Write(data)
	if err != nil {
		lobw.txCommit = false
	}
	return n, err
}

func (lobw *lobWriteWrapper) Close() error {
	defer lobw.finisher(lobw.txCommit)

	err := lobw.lob.Close()
	if err != nil {
		lobw.txCommit = false
		return err
	}
	return nil
}
