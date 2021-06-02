package global

import "github.com/prometheus/client_golang/prometheus"

const (
	// Exporter Namespace.
	Namespace = "my_one_exporter"
)

func NewDesc(subsystem, name, help string) *prometheus.Desc {
	return prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subsystem, name),
		help, nil, nil,
	)
}
