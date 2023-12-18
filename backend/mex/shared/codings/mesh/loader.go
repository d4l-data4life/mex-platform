package mesh

import (
	"compress/zlib"
	"context"
	"io"

	"google.golang.org/protobuf/types/known/anypb"

	"github.com/d4l-data4life/mex/mex/shared/blobs"
	"github.com/d4l-data4life/mex/mex/shared/codings"
	"github.com/d4l-data4life/mex/mex/shared/codings/csrepo"
	"github.com/d4l-data4life/mex/mex/shared/utils/async"
)

func NewBlobStoreLoader(blobStore blobs.BlobStore) csrepo.CodingsetLoader {
	return func(config *anypb.Any) async.Promise[codings.Codingset] {
		return async.New(func(resolve async.Resolver[codings.Codingset], reject async.Rejecter) {
			var cfg codings.BlobStoreCodingsetSourceConfig
			err := config.UnmarshalTo(&cfg)
			if err != nil {
				reject(err)
				return
			}

			r, err := blobStore.GetReadCloser(context.Background(), cfg.BlobName, cfg.BlobType)
			if err != nil {
				reject(err)
				return
			}
			defer r.Close()

			rc, err := zlib.NewReader(r)
			if err != nil {
				reject(err)
				return
			}
			buf, err := io.ReadAll(rc)
			if err != nil {
				reject(err)
				return
			}

			codingset, err := NewCodingsetBytes(buf)
			if err != nil {
				reject(err)
				return
			}

			resolve(codingset)
		})
	}
}
