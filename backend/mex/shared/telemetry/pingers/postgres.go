package pingers

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/d4l-data4life/mex/mex/shared/utils"

	"github.com/d4l-data4life/mex/mex/services/metadata/business/datamodel"
)

type PostgresPinger struct {
	LastStoredItemsCount int64

	lastError error
	quit      chan<- struct{}
}

func (p *PostgresPinger) LastError() error {
	return p.lastError
}

func (p *PostgresPinger) Stop() {
	close(p.quit)
}

func NewPostgresPinger(client *pgxpool.Pool, updateInterval time.Duration) *PostgresPinger {
	if client == nil {
		panic("Postgres client is nil")
	}

	p := PostgresPinger{
		lastError: fmt.Errorf("Postgres pinger did not run yet"),
	}

	p.quit = utils.ExponentialTicker(time.Second, updateInterval, func() {
		queries := datamodel.New(client)
		count, err := queries.DbGetItemsCount(context.Background())
		p.lastError = err
		if err == nil {
			p.LastStoredItemsCount = count
		} else {
			p.LastStoredItemsCount = -1
		}
	})

	return &p
}

func (p *PostgresPinger) RegisterMetrics(namespace string, reg *prometheus.Registry) {
	m0 := prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "stored_items_count",
		Help:      "Number of items in Postgres",
	}, func() float64 {
		return float64(p.LastStoredItemsCount)
	})
	reg.MustRegister(m0)
}
