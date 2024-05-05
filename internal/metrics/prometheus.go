package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sbnarra/bckupr/internal/utils/errors"
)

var (
	backupDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "backup_duration_seconds",
		Help: "Backup duration in seconds",
	}, []string{"id", "action"})

	backupSuccessTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "backup_success_total",
		Help: "Number of successful backups",
	}, []string{"id", "action"})

	backupErrorTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "backup_error_total",
		Help: "Number of failed backups",
	}, []string{"id", "action"})

	restoreDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "restore_duration_seconds",
		Help: "Restores duration in seconds",
	}, []string{"id", "action"})

	restoreSuccessTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "restore_success_total",
		Help: "Number of successful restores",
	}, []string{"id", "action"})

	restoreErrorTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "restore_error_total",
		Help: "Number of failed restores",
	}, []string{"id", "action"})
)

type Metrics struct {
	id        string
	action    string
	volume    string
	jobLabels []string
	timer     *prometheus.Timer

	successTotal *prometheus.CounterVec
	errorTotal   *prometheus.CounterVec
}

func Backup(id string, volume string) Metrics {
	jobLabels := []string{id, "backup"}
	timer := prometheus.NewTimer(backupDuration.WithLabelValues(jobLabels...))
	return Metrics{
		id:        id,
		action:    "backup",
		jobLabels: jobLabels,
		volume:    volume,
		timer:     timer,

		successTotal: backupSuccessTotal,
		errorTotal:   backupErrorTotal,
	}
}

func Restore(id string, volume string) Metrics {
	jobLabels := []string{id, "restore"}
	timer := prometheus.NewTimer(restoreDuration.WithLabelValues(jobLabels...))
	return Metrics{
		id:        id,
		action:    "restore",
		jobLabels: jobLabels,
		volume:    volume,
		timer:     timer,

		successTotal: restoreSuccessTotal,
		errorTotal:   restoreErrorTotal,
	}
}

func (m Metrics) OnComplete(err *errors.Error) {
	m.timer.ObserveDuration()
	if err != nil {
		m.errorTotal.WithLabelValues(m.jobLabels...).Inc()
	} else {
		m.successTotal.WithLabelValues(m.jobLabels...).Inc()
	}
}
