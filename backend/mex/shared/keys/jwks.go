package keys

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"

	L "github.com/d4l-data4life/mex/mex/shared/log"
)

type TokenValidator interface {
	ValidateJWT(token string) (*jwt.Token, error)
	IsReady() bool
}

type RemoteKeyStore struct {
	keyStore *jwk.AutoRefresh
	uri      string
	ready    bool
}

type RemoteKeyStoreOptions struct {
	Log                         L.Logger
	URI                         string
	TLSConfig                   *tls.Config
	InternalAuthServiceHostname string
}

func NewRemoteKeyStore(ctx context.Context, options RemoteKeyStoreOptions) (*RemoteKeyStore, error) {
	u, err := url.Parse(options.URI)
	if err != nil {
		return nil, err
	}

	autoRefresher := jwk.NewAutoRefresh(ctx)

	// We require the protocol for querying the keys to be https unless it is our internal auth service.
	if u.Scheme != "https" {
		if u.Hostname() != options.InternalAuthServiceHostname {
			return nil, fmt.Errorf("key store URL can only be HTTP if used for internal auth service")
		}
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: options.TLSConfig,
		},
	}
	autoRefresher.Configure(options.URI, jwk.WithHTTPClient(client))

	set, err := autoRefresher.Refresh(ctx, options.URI)
	if err != nil {
		return nil, err
	}

	options.Log.Info(ctx, L.Messagef("# of retrieved keys: %d", set.Len()), L.Phase("jwks"))

	return &RemoteKeyStore{
		keyStore: autoRefresher,
		uri:      options.URI,
		ready:    true,
	}, nil
}

func (ks *RemoteKeyStore) IsReady() bool {
	return ks.ready
}

func (ks *RemoteKeyStore) ValidateJWT(token string) (*jwt.Token, error) {
	var t jwt.Token

	if ks.keyStore == nil {
		return nil, errors.New("no remote key registry configured")
	}

	// Fetch will honor all HTTP cache headers that may be sent by the keys endpoint.
	// That is, we do not do an HTTP request each time!
	set, err := ks.keyStore.Fetch(context.Background(), ks.uri)
	if err != nil {
		return nil, err
	}

	t, err = jwt.Parse([]byte(token), jwt.WithValidate(true), jwt.InferAlgorithmFromKey(true), jwt.WithKeySet(set))
	if err == nil {
		return &t, nil
	}

	return nil, err
}
