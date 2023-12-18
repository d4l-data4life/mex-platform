package pingers

import (
	"context"
	"fmt"
	"time"

	"github.com/d4l-data4life/mex/mex/shared/solr"
	"github.com/d4l-data4life/mex/mex/shared/utils"
)

type SolrPinger struct {
	lastError error
	quit      chan<- struct{}
}

func (p *SolrPinger) LastError() error {
	return p.lastError
}

func (p *SolrPinger) Stop() {
	close(p.quit)
}

func NewSolrPinger(client solr.ClientAPI, updateInterval time.Duration) *SolrPinger {
	if client == nil {
		panic("Solr client is nil")
	}

	p := SolrPinger{
		lastError: fmt.Errorf("Solr pinger did not run yet"),
	}

	p.quit = utils.ExponentialTicker(time.Second, updateInterval, func() {
		p.lastError = client.Ping(context.Background())
	})

	return &p
}
