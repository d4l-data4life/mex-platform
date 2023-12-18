package certs

import (
	"context"
	"crypto/x509"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	L "github.com/d4l-data4life/mex/mex/shared/log"
	"github.com/d4l-data4life/mex/mex/shared/try"
)

type LoadParams struct {
	AdditionalCACertsFiles []string
	AdditionalCACertsPEM   []byte

	ConnectionAttempts uint32
	ConnectionPause    time.Duration

	Log L.Logger
}

func LoadCerts(ctx context.Context, params LoadParams) (*x509.CertPool, error) {
	params.Log.Info(ctx, L.Messagef("additional CA PEM files: %v", params.AdditionalCACertsFiles))

	caPEMReaders := []io.Reader{}
	for _, pemFile := range params.AdditionalCACertsFiles {
		trimmedFilename := strings.TrimSpace(pemFile)
		if trimmedFilename == "" {
			continue
		}

		r, err := try.Try(ctx, try.Task[io.Reader]{
			Desc:          "load PEM file: " + trimmedFilename,
			Phase:         "startup",
			PauseStrategy: try.NewMaxAttemptsConstantPauseStrategy(params.ConnectionAttempts, params.ConnectionPause),
			Log:           params.Log,
			Func: func() (io.Reader, error) {
				return os.Open(trimmedFilename)
			},
		})
		if err != nil {
			return nil, fmt.Errorf("could not load PEM file %s: %w", trimmedFilename, err)
		}
		caPEMReaders = append(caPEMReaders, r)
	}
	if len(params.AdditionalCACertsPEM) > 0 {
		params.Log.Info(ctx, L.Messagef("PEM content in env var, length: %d", len(params.AdditionalCACertsPEM)))
		if params.AdditionalCACertsPEM[0] != ' ' {
			caPEMReaders = append(caPEMReaders, strings.NewReader(string(params.AdditionalCACertsPEM)))
		} else {
			params.Log.Info(ctx, L.Message("ignoring empty PEM parameter"))
		}
	}

	rootCAs, err := GetTLSClientConfig(ctx, params.Log, caPEMReaders)
	if err != nil {
		return nil, err
	}

	return rootCAs, nil
}
