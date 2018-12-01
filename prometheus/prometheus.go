package prometheus

import (
	"github.com/krpn/youtrack-issues-prometheus-exporter/model"
	pr "github.com/prometheus/client_golang/prometheus"
)

// Metrics describes Prometheus metric collector.
type Metrics struct {
	issues gaugeIniter
	errors counterIniter
}

// New creates Metrics.
func New() *Metrics {
	issues := pr.NewGaugeVec(
		pr.GaugeOpts{
			Subsystem: "youtrack",
			Name:      "issues",
			Help:      "Query issues",
		},
		[]string{"query", "id", "title"},
	)

	errors := pr.NewCounterVec(
		pr.CounterOpts{
			Subsystem: "youtrack",
			Name:      "errors",
			Help:      "Errors counter",
		},
		[]string{"query", "error"},
	)

	pr.MustRegister(issues)
	pr.MustRegister(errors)

	return &Metrics{
		issues: issues,
		errors: errors,
	}
}

// EnableMonitoring turns on metric for issue.
func (p *Metrics) EnableMonitoring(queryName string, issue model.Issue) {
	p.issues.WithLabelValues(queryName, issue.ID, issue.Title).Set(1)
}

// DisableMonitoring turns off metric for issue.
func (p *Metrics) DisableMonitoring(queryName string, issue model.Issue) {
	p.issues.WithLabelValues(queryName, issue.ID, issue.Title).Set(0)
}

// ErrorInc increments metric for error
func (p *Metrics) ErrorInc(queryName string, err error) {
	p.errors.WithLabelValues(queryName, err.Error()).Inc()
}

//go:generate mockgen -destination=prometheus_metrics_mocks.go -package=prometheus github.com/prometheus/client_golang/prometheus Counter,Gauge
//go:generate mockgen -source=prometheus.go -destination=prometheus_mocks.go -package=prometheus doc github.com/golang/mock/gomock

type counterIniter interface {
	WithLabelValues(lvs ...string) pr.Counter
}

type gaugeIniter interface {
	WithLabelValues(lvs ...string) pr.Gauge
}
