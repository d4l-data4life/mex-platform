package rdb

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"

	"github.com/go-redis/redis/v8"

	L "github.com/d4l-data4life/mex/mex/shared/log"
	"github.com/d4l-data4life/mex/mex/shared/try"
)

type InitRedisSupport struct {
	RootCAs *x509.CertPool

	Hostname string
	Port     uint32
	Password string
	DB       uint32
	UseTLS   bool

	RetryStrategy try.PauseStrategy
}

func (params *InitRedisSupport) Init(ctx context.Context, log L.Logger) (*redis.Client, error) {
	var redisTLSConfig *tls.Config
	if params.UseTLS {
		//nolint:gosec // gosec complains about this issue: https://github.com/go-redis/redis/issues/1553
		redisTLSConfig = &tls.Config{
			InsecureSkipVerify: false,
			ServerName:         params.Hostname,
			RootCAs:            params.RootCAs,
		}
	}
	rdb, err := try.Try(ctx, try.Task[*redis.Client]{
		Desc:          "Redis connection",
		Phase:         "startup",
		PauseStrategy: params.RetryStrategy,
		Log:           log,
		Func: func() (*redis.Client, error) {
			rdb := redis.NewClient(&redis.Options{
				Addr:      fmt.Sprintf("%s:%d", params.Hostname, params.Port),
				Password:  params.Password,
				DB:        int(params.DB),
				TLSConfig: redisTLSConfig, // nil is fine as it means no TLS
			})

			err := rdb.Ping(ctx).Err()
			if err != nil {
				return nil, err
			}

			return rdb, nil
		},
	})
	if err != nil {
		return nil, err
	}

	return rdb, nil
}
