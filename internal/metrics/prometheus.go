package metrics

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	actionDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "action_duration_seconds",
		Help: "Action duration in seconds",
	}, []string{"id", "action"})

	actionSuccessTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "action_success_total",
		Help: "Number of successful X",
	}, []string{"id", "action"})

	actionErrorTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "action_error_total",
		Help: "Number of failed X",
	}, []string{"id", "action"})
)

type Metrics struct {
	id        string
	action    string
	volume    string
	jobLabels []string
	timer     *prometheus.Timer
}

func New(id string, action string, volume string) Metrics {
	jobLabels := []string{id, action}
	timer := prometheus.NewTimer(actionDuration.WithLabelValues(jobLabels...))
	return Metrics{
		id:        id,
		action:    action,
		jobLabels: jobLabels,
		volume:    volume,
		timer:     timer,
	}
}

func (m Metrics) OnComplete(err error) {
	m.timer.ObserveDuration()

	if err != nil {
		actionErrorTotal.WithLabelValues(m.jobLabels...).Inc()
		fmt.Printf("Error: %s (id=%s, volume=%s, err=%s)", m.action, m.id, m.volume, err)
	} else {
		actionSuccessTotal.WithLabelValues(m.jobLabels...).Inc()
		fmt.Printf("Success: %s (id=%s, volume=%s)", m.action, m.id, m.volume)
	}
}
