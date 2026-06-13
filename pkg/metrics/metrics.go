package metrics

import "github.com/prometheus/client_golang/prometheus"

type Metrics struct {
	QueueDepth    *prometheus.GaugeVec
	JobsSubmitted *prometheus.CounterVec
}

func NewMetrics(reg prometheus.Registerer) *Metrics {
	m := &Metrics{
		JobsSubmitted: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "inference_jobs_submitted_total",
				Help: "Total number of jobs submitted",
			},
			[]string{"type"}, // label — lets you filter by classification vs regression
		),

		QueueDepth: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "inference_queue_depth",
				Help: "Current number of jobs waiting in queue",
			},
			[]string{"type"},
		),
	}
	reg.MustRegister(m.JobsSubmitted)
	reg.MustRegister(m.QueueDepth)
	return m
}
