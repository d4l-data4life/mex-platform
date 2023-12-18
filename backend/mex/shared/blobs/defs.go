package blobs

import (
	"context"
	"io"
)

type BlobStore interface {
	List(ctx context.Context) ([]*BlobInfo, error)
	Delete(ctx context.Context, blobName string, blobType string) error

	GetReadCloser(ctx context.Context, blobName string, blobType string) (io.ReadCloser, error)
	GetWriteCloser(ctx context.Context, blobName string, blobType string, append bool) (io.WriteCloser, error)
}
