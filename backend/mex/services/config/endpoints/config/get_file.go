package config

import (
	"bytes"
	"context"
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"strings"

	"google.golang.org/grpc/codes"

	E "github.com/d4l-data4life/mex/mex/shared/errstat"
	L "github.com/d4l-data4life/mex/mex/shared/log"

	pbConfig "github.com/d4l-data4life/mex/mex/services/config/endpoints/config/pb"
)

func (svc *Service) GetFile(ctx context.Context, request *pbConfig.GetFileRequest) (*pbConfig.GetFileResponse, error) {
	svc.Log.Info(ctx, L.Messagef("name: '%s' (env path: '%s')", request.Name, svc.EnvPath))

	if svc.fs == nil {
		return nil, E.MakeGRPCStatus(codes.Internal, "nothing checked out yet").Err()
	}

	svc.mu.RLock()
	defer svc.mu.RUnlock()

	effectiveFileName := "./" + svc.EnvPath + "/" + request.Name
	info, err := svc.fs.Stat(effectiveFileName)
	if err != nil {
		if err != os.ErrNotExist {
			return nil, err
		}

		// Try again lowercase
		effectiveFileName = strings.ToLower(effectiveFileName)
		info, err = svc.fs.Stat(effectiveFileName)
		if err != nil {
			if err == os.ErrNotExist {
				return nil, E.MakeGRPCStatus(codes.NotFound, "file not found", E.Cause(err), E.DevMessagef("file not found: %s", effectiveFileName)).Err()
			}
			return nil, err
		}
	}

	if info.IsDir() {
		effectiveFileName = fmt.Sprintf("%s/index.json", effectiveFileName)
		_, err = svc.fs.Stat(effectiveFileName)
		if err != nil {
			if err == os.ErrNotExist {
				return nil, E.MakeGRPCStatus(codes.NotFound, "file not found", E.Cause(err), E.DevMessagef("file not found: %s", effectiveFileName)).Err()
			}
			return nil, err
		}
	}

	svc.Log.Info(ctx, L.Messagef("effective name: %s", effectiveFileName))

	f, err := svc.fs.Open(effectiveFileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf := bytes.Buffer{}
	n, err := buf.ReadFrom(f)
	if err != nil {
		return nil, err
	}
	svc.Log.Info(ctx, L.Messagef("bytes read: %d", n))

	return &pbConfig.GetFileResponse{
		MimeType: mime.TypeByExtension(filepath.Ext(effectiveFileName)),
		Content:  buf.Bytes(),
	}, nil
}
