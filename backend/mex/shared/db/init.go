package db

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"

	L "github.com/d4l-data4life/mex/mex/shared/log"
	"github.com/d4l-data4life/mex/mex/shared/try"
)

type InitDatabaseSupport struct {
	RootCAs *x509.CertPool

	User       string
	Password   string
	Hostname   string
	Port       uint32
	Name       string
	SearchPath []string
	SSLMode    string

	RetryStrategy try.PauseStrategy
}

func (params *InitDatabaseSupport) Init(ctx context.Context, log L.Logger) (*pgxpool.Pool, error) {
	// This is the database connection object which we will pass to all the services below.
	poolConfig, err := pgxpool.ParseConfig(fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		params.User, params.Password, params.Hostname, params.Port, params.Name, params.SSLMode))
	if err != nil {
		return nil, err
	}
	if params.SSLMode != "disable" {
		//nolint:gosec // gosec complains about this issue: https://github.com/go-redis/redis/issues/1553
		poolConfig.ConnConfig.TLSConfig = &tls.Config{
			InsecureSkipVerify: false,
			ServerName:         params.Hostname,
			RootCAs:            params.RootCAs,
		}
	}

	// This map is initialized by the ParseConfig function, so it is safe to assign.
	poolConfig.ConnConfig.RuntimeParams["search_path"] = strings.Join(params.SearchPath, ",")

	// Establish the connection; try multiple times in case the database is not online yet.
	// (can happen `during docker-compose up`)
	return try.Try(ctx, try.Task[*pgxpool.Pool]{
		Desc:          "database connection",
		Phase:         "startup",
		PauseStrategy: params.RetryStrategy,
		Log:           log,
		Func: func() (*pgxpool.Pool, error) {
			return pgxpool.NewWithConfig(ctx, poolConfig)
		},
	})
}
