package config

import (
	"context"

	"github.com/go-git/go-billy/v5"

	pbConfig "github.com/d4l-data4life/mex/mex/services/config/endpoints/config/pb"
)

func (svc *Service) ListConfig(ctx context.Context, request *pbConfig.ListConfigRequest) (*pbConfig.ListConfigResponse, error) {
	svc.mu.RLock()
	defer svc.mu.RUnlock()

	fileNames := []string{}

	err := listFiles(svc.fs, &fileNames, ".")
	if err != nil {
		return nil, err
	}

	return &pbConfig.ListConfigResponse{
		FileName: fileNames,
	}, nil
}

func listFiles(fs billy.Filesystem, files *[]string, f string) error {
	info, err := fs.Stat(f)
	if err != nil {
		return err
	}

	if !info.IsDir() {
		*files = append(*files, f)
		return nil
	}

	infos, err := fs.ReadDir(f)
	if err != nil {
		return err
	}

	for _, info := range infos {
		err = listFiles(fs, files, f+"/"+info.Name())
		if err != nil {
			return err
		}
	}

	return nil
}
