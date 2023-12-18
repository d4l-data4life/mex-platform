package blobs

import (
	"compress/zlib"
	"context"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc/codes"

	sharedBlobs "github.com/d4l-data4life/mex/mex/shared/blobs"
	"github.com/d4l-data4life/mex/mex/shared/blobs/pglo"
	"github.com/d4l-data4life/mex/mex/shared/codings"
	"github.com/d4l-data4life/mex/mex/shared/codings/mesh"
	E "github.com/d4l-data4life/mex/mex/shared/errstat"
	"github.com/d4l-data4life/mex/mex/shared/hints"
	L "github.com/d4l-data4life/mex/mex/shared/log"

	pbBlobs "github.com/d4l-data4life/mex/mex/services/metadata/endpoints/blobs/pb"
)

type Service struct {
	Log L.Logger

	DB              *pgxpool.Pool
	MasterTableName string

	pbBlobs.UnimplementedBlobsServer
}

func (svc *Service) CreateBlob(ctx context.Context, request *pbBlobs.CreateBlobRequest) (*pbBlobs.CreateBlobResponse, error) {
	svc.Log.Info(ctx, L.Messagef("request data size: %d", len(request.Data)))

	if strings.TrimSpace(request.BlobName) == "" {
		return nil, E.MakeGRPCStatus(codes.InvalidArgument, "blob name is empty").Err()
	}

	if strings.TrimSpace(request.BlobType) == "" {
		return nil, E.MakeGRPCStatus(codes.InvalidArgument, "blob type is empty").Err()
	}

	blobStore := pglo.PostgresLargeObjectStore{
		DB:              svc.DB,
		MasterTableName: svc.MasterTableName,
		Log:             svc.Log,
	}

	w, err := blobStore.GetWriteCloser(ctx, request.BlobName, request.BlobType, request.Append)
	if err != nil {
		return nil, E.MakeGRPCStatus(codes.Internal, err.Error()).Err()
	}
	defer w.Close()

	n, err := w.Write(request.Data)
	if err != nil {
		return nil, E.MakeGRPCStatus(codes.InvalidArgument, err.Error()).Err()
	}

	hints.HintHTTPStatusCode(ctx, http.StatusCreated)
	return &pbBlobs.CreateBlobResponse{BytesWritten: int32(n)}, nil
}

func (svc *Service) ListBlobs(ctx context.Context, request *pbBlobs.ListBlobsRequest) (*pbBlobs.ListBlobsResponse, error) {
	blobStore := pglo.PostgresLargeObjectStore{
		DB:              svc.DB,
		MasterTableName: svc.MasterTableName,
		Log:             svc.Log,
	}

	blobInfos, err := blobStore.List(ctx)
	if err != nil {
		return nil, E.MakeGRPCStatus(codes.Internal, err.Error()).Err()
	}

	response := pbBlobs.ListBlobsResponse{BlobInfos: make([]*sharedBlobs.BlobInfo, len(blobInfos))}
	for i := range blobInfos {
		response.BlobInfos[i] = &sharedBlobs.BlobInfo{
			BlobName: blobInfos[i].BlobName,
			BlobType: blobInfos[i].BlobType,
		}
	}
	return &response, nil
}

func (svc *Service) DeleteBlob(ctx context.Context, request *pbBlobs.DeleteBlobRequest) (*pbBlobs.DeleteBlobResponse, error) {
	svc.Log.Info(ctx, L.Messagef("delete blob: %s", request.BlobName))

	blobStore := pglo.PostgresLargeObjectStore{
		DB:              svc.DB,
		MasterTableName: svc.MasterTableName,
		Log:             svc.Log,
	}

	err := blobStore.Delete(ctx, request.BlobName, request.BlobType)
	if err != nil {
		return nil, E.MakeGRPCStatus(E.CodeFrom(err), err.Error()).Err()
	}

	hints.HintHTTPStatusCode(ctx, http.StatusNoContent)
	return &pbBlobs.DeleteBlobResponse{}, nil
}

func (svc *Service) GetBlob(ctx context.Context, request *pbBlobs.GetBlobRequest) (*pbBlobs.GetBlobResponse, error) {
	svc.Log.Info(ctx, L.Messagef("get blob: %s / %s", request.BlobName, request.BlobType))

	blobStore := pglo.PostgresLargeObjectStore{
		DB:              svc.DB,
		MasterTableName: svc.MasterTableName,
		Log:             svc.Log,
	}

	r, err := blobStore.GetReadCloser(ctx, request.BlobName, request.BlobType)
	if err != nil {
		return nil, E.MakeGRPCStatus(E.CodeFrom(err), err.Error()).Err()
	}
	defer r.Close()

	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return &pbBlobs.GetBlobResponse{Data: data}, nil
}

func (svc *Service) MeshTest(ctx context.Context, request *pbBlobs.MeshTestRequest) (*pbBlobs.MeshTestResponse, error) {
	blobStore := pglo.PostgresLargeObjectStore{
		DB:              svc.DB,
		MasterTableName: svc.MasterTableName,
		Log:             svc.Log,
	}

	r, err := blobStore.GetReadCloser(ctx, request.BlobName, request.BlobType)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	rc, err := zlib.NewReader(r)
	if err != nil {
		return nil, err
	}
	buf, err := io.ReadAll(rc)
	if err != nil {
		return nil, err
	}
	svc.Log.Info(ctx, L.Messagef("inflated buffer size: %d", len(buf)))

	var codingset codings.Codingset

	switch request.LoadingMode {
	case pbBlobs.MeshTestRequest_LOADING_MODE_IN_MEMORY:
		codingset, err = mesh.NewCodingsetBytes(buf)
		if err != nil {
			return nil, err
		}
	case pbBlobs.MeshTestRequest_LOADING_MODE_TEMP_FILE:
		f, err := os.CreateTemp(".", "mesh-*")
		if err != nil {
			return nil, err
		}
		_, _ = f.Write(buf)
		_ = f.Close()
		defer os.Remove(f.Name()) // delete temp file after leaving this function

		codingset, err = mesh.NewCodingsetFile(f)
		if err != nil {
			return nil, err
		}
	}

	if request.RunGc {
		runGC()
	}

	defer codingset.Close()

	start := time.Now()
	mh, err := codingset.GetMainHeadings()
	if err != nil {
		svc.Log.Warn(ctx, L.Message(err.Error()))
	}
	svc.Log.Info(ctx, L.Messagef("duration: %s", time.Since(start)))
	svc.Log.Info(ctx, L.Messagef("# of main headings: %d", len(mh)))

	if len(mh) == 0 {
		return nil, E.MakeGRPCStatus(codes.Internal, "no main headings found").Err()
	}

	if request.BagSize <= 0 {
		return nil, E.MakeGRPCStatus(codes.InvalidArgument, "bag size must be positive").Err()
	}
	if request.Iterations <= 0 {
		return nil, E.MakeGRPCStatus(codes.InvalidArgument, "iterations must be positive").Err()
	}

	headings := make([]string, request.BagSize)
	for n := 0; n < int(request.Iterations); n++ {
		for h := range headings {
			//nolint:gosec
			headings[h] = mh[rand.Intn(len(mh))]
		}

		start = time.Now()
		terms, _ := codingset.ResolveMainHeadings(headings, "en")
		svc.Log.Info(ctx, L.Messagef("duration: EN  : %s (%d)", time.Since(start), len(terms)))
		if request.ShowTerms {
			fmt.Printf("%v\n", terms)
		}

		start = time.Now()
		terms, _ = codingset.ResolveMainHeadings(headings, "de")
		svc.Log.Info(ctx, L.Messagef("duration: DE  : %s (%d)", time.Since(start), len(terms)))
		if request.ShowTerms {
			fmt.Printf("%v\n", terms)
		}

		start = time.Now()
		tnums, _ := codingset.ResolveTreeNumbers(headings)
		svc.Log.Info(ctx, L.Messagef("duration: tree: %s (%d)", time.Since(start), len(tnums)))
		if request.ShowTerms {
			fmt.Printf("%v\n", tnums)
		}
	}

	info, err := codingset.Info()
	if err != nil {
		svc.Log.Warn(ctx, L.Message(err.Error()))
	}

	return &pbBlobs.MeshTestResponse{
		Info:          info,
		DistinctCount: int32(codingset.Count()),
	}, nil
}

func runGC() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	fmt.Println("GC - BEFORE:")
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)

	start := time.Now()
	runtime.GC()
	fmt.Printf("GC - AFTER: duration: %s\n", time.Since(start))

	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

//nolint:gomnd
func bToMb(b uint64) uint64 {
	return b / (1 << 20)
}
