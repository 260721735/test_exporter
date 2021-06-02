package collector

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"sync"
	"test_exporter/global"
	"test_exporter/scrape"
	"time"
)

const (
	// Subsystem(s).
	exporter = "exporter"
)

var (
	scrapeDurationDesc = prometheus.NewDesc(
		prometheus.BuildFQName(global.Namespace, exporter, "collector_duration_seconds"),
		"Collector time duration.",
		[]string{"collector"}, nil,
	)
)

type Metrics struct {
	ExporterUp prometheus.Gauge
}

func NewMetrics() Metrics {
	return Metrics{
		ExporterUp: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: global.Namespace,
			Name:      "up",
			Help:      "Whether the datacenter is up.",
		}),
	}
}

type Exporter struct {
	ctx      context.Context
	dsn      string
	scrapers []scrape.Scraper
	metrics  Metrics
}

func New(ctx context.Context, dsn string, metrics Metrics, scrapers []scrape.Scraper) *Exporter {
	return &Exporter{
		ctx:      ctx,
		dsn:      dsn,
		scrapers: scrapers,
		metrics:  metrics,
	}
}
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.metrics.ExporterUp.Desc()
}
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.scrape(e.ctx, ch)
	ch <- e.metrics.ExporterUp
}
func (e *Exporter) scrape(ctx context.Context, ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(scrapeDurationDesc, prometheus.GaugeValue, 0.01, "version")
	e.metrics.ExporterUp.Set(1)
	var wg sync.WaitGroup
	defer wg.Wait()
	for _, scraper := range e.scrapers {
		wg.Add(1)
		go func(scraper scrape.Scraper) {
			defer wg.Done()
			label := "collect." + scraper.Name()
			scrapeTime := time.Now()
			if err := scraper.Scrape(ctx, "dc", ch); err != nil {
				log.Println(err)
			}
			ch <- prometheus.MustNewConstMetric(scrapeDurationDesc, prometheus.GaugeValue, time.Since(scrapeTime).Seconds(), label)
		}(scraper)
	}
}
