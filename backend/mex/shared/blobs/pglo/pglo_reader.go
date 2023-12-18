package pglo

import (
	"context"
	"fmt"
	"io"

	"github.com/jackc/pgx/v5"

	"github.com/d4l-data4life/mex/mex/shared/blobs"
	L "github.com/d4l-data4life/mex/mex/shared/log"
)

func (store *PostgresLargeObjectStore) List(ctx context.Context) ([]*blobs.BlobInfo, error) {
	if store.DB == nil {
		return nil, fmt.Errorf("store DB is nil")
	}

	rows, err := store.DB.Query(ctx, fmt.Sprintf("SELECT blob_name, blob_type FROM %s ORDER BY blob_name ASC, blob_type ASC", store.MasterTableName))
	if err != nil {
		return nil, fmt.Errorf("error reading blob master table: %w", err)
	}
	defer rows.Close()

	var blobInfos []*blobs.BlobInfo
	for rows.Next() {
		var bi blobs.BlobInfo
		if err := rows.Scan(&bi.BlobName, &bi.BlobType); err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		blobInfos = append(blobInfos, &bi)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return blobInfos, nil
}

type lobReadWrapper struct {
	lob      *pgx.LargeObject
	finisher func()
}

func (store *PostgresLargeObjectStore) GetReadCloser(ctx context.Context, blobName string, blobType string) (io.ReadCloser, error) {
	if store.DB == nil {
		return nil, fmt.Errorf("store DB is nil")
	}

	tx, err := store.DB.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("error beginning transaction: %w", err)
	}

	row := tx.QueryRow(ctx, fmt.Sprintf("SELECT blob_oid FROM %s WHERE blob_name = $1 AND blob_type = $2", store.MasterTableName), blobName, blobType)
	if row == nil {
		return nil, fmt.Errorf("blob not found: %s/%s", blobName, blobType)
	}

	var oid uint32
	err = row.Scan(&oid)
	if err != nil {
		return nil, fmt.Errorf("invalid blob oid: %w", err)
	}
	store.logInfo(ctx, L.Messagef("OID: %d", oid))

	lo := tx.LargeObjects()
	lob, err := lo.Open(ctx, oid, pgx.LargeObjectModeRead)
	if err != nil {
		return nil, fmt.Errorf("error reading blob '%s/%s': %w", blobName, blobType, err)
	}

	return &lobReadWrapper{
		lob: lob,
		finisher: func() {
			ok := true
			finishTx(ctx, tx, store.Log, &ok)
		},
	}, nil
}

func (lobr *lobReadWrapper) Read(p []byte) (int, error) {
	return lobr.lob.Read(p)
}

func (lobr *lobReadWrapper) Close() error {
	_ = lobr.lob.Close()
	lobr.finisher()
	return nil
}
