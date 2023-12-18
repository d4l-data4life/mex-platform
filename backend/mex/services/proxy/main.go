package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/d4l-data4life/mex/mex/shared/cfg"
	L "github.com/d4l-data4life/mex/mex/shared/log"
	"github.com/d4l-data4life/mex/mex/shared/log/emit"

	"github.com/d4l-data4life/mex/mex/services/proxy/proxycfg"
)

const pathReadinessProbe = "/probes/readiness"

func makeProxyHandler(log L.Logger, originMetadata, originQuery, originIndex, originAuth, originCms, originConfig, originWebapp *url.URL) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.LogHTTPRequest(r)

		var u *url.URL

		switch {
		case strings.HasPrefix(r.RequestURI, "/api/v0/query"):
			u = originQuery
		case strings.HasPrefix(r.RequestURI, "/api/v0/oauth"):
			u = originAuth
		case strings.HasPrefix(r.RequestURI, "/api/v0/metadata/index"):
			u = originIndex
		case strings.HasPrefix(r.RequestURI, "/api/v0/config"):
			u = originConfig
		case strings.HasPrefix(r.RequestURI, "/cms"):
			u = originCms
		case strings.HasPrefix(r.RequestURI, "/api"):
			u = originMetadata
		case strings.HasPrefix(r.RequestURI, "/api/v1/events"):
			w.WriteHeader(http.StatusOK)
			return
		case strings.HasPrefix(r.RequestURI, "/probes/metadata/readiness"):
			u = originMetadata
			r.URL.Path = pathReadinessProbe
		case strings.HasPrefix(r.RequestURI, "/probes/query/readiness"):
			u = originQuery
			r.URL.Path = pathReadinessProbe
		case strings.HasPrefix(r.RequestURI, "/probes/index/readiness"):
			u = originIndex
			r.URL.Path = pathReadinessProbe
		case strings.HasPrefix(r.RequestURI, "/probes/auth/readiness"):
			u = originAuth
			r.URL.Path = pathReadinessProbe
		case strings.HasPrefix(r.RequestURI, "/probes/config/readiness"):
			u = originConfig
			r.URL.Path = pathReadinessProbe
		default: // go to webapp
			u = originWebapp
		}

		r.Close = true
		proxy := httputil.NewSingleHostReverseProxy(u)

		nw := L.NewPeekingResponseWriter(w)
		proxy.ServeHTTP(nw, r)

		log.LogHTTPResponse(r, nw.StatusCode, nw.BodyCount, start, true)
	}
}

func main() {
	log, err := L.New("MEx Proxy", "develop", &emit.WriterEmitter{Writer: os.Stdout})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := run(log); err != nil {
		log.Error(context.TODO(), L.Message(err.Error()))
		os.Exit(1)
	}
}

func run(log L.Logger) error {
	ctx := context.Background()

	config := proxycfg.ProxyConfig{}
	dump, err := cfg.InitConfig(log, &cfg.OSEnvs{}, "MEX", "proxy", &config)
	if err != nil {
		log.Error(ctx, L.Messagef("could not determine config: %s", err.Error()))
		os.Exit(2)
	}
	log.Info(ctx, L.Messagef("effective config:\n%s", dump))

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	serverErrors := make(chan error, 1)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	s := http.Server{
		Addr:              config.Addr,
		ReadTimeout:       config.Timeout.AsDuration(),
		ReadHeaderTimeout: config.Timeout.AsDuration(),
		WriteTimeout:      config.Timeout.AsDuration(),
		IdleTimeout:       config.Timeout.AsDuration(),
		Handler: makeProxyHandler(log,
			URLMustParse(config.Origin.Metadata),
			URLMustParse(config.Origin.Query),
			URLMustParse(config.Origin.Index),
			URLMustParse(config.Origin.Auth),
			URLMustParse(config.Origin.Cms),
			URLMustParse(config.Origin.Config),
			URLMustParse(config.Origin.Webapp),
		),
	}

	go func() {
		log.Info(context.TODO(), L.Messagef("proxy started: %s (PID %d)", config.Addr, os.Getpid()))
		serverErrors <- s.ListenAndServe()
	}()

	// ------------------------------------------------------------------------------------------------------
	// Shutdown logic

	// Block main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case sig := <-shutdown:
		log.Info(ctx, L.Messagef("shutdown initiated, signal: %v", sig))
		return nil
	}
}

func URLMustParse(rawURL string) *url.URL {
	u, err := url.Parse(rawURL)
	if err != nil {
		panic(err)
	}
	return u
}
