package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

func init() {
	RegisterExporter("shovel", newExporterShovel)
}

var (
	shovelLabels     = []string{"dest_uri", "name", "node", "state", "timestamp", "type"}
	shovelLabelsKeys = []string{"dest_uri", "name", "node", "state", "timestamp", "type"}

	shovelGaugeVec = map[string]*prometheus.GaugeVec{
		"state":     newGaugeVec("shovel_state", "shovel state", shovelLabels),
		"timestamp": newGaugeVec("shovel_timestamp", "shovel timestamp", shovelLabels),
	}
)

type exporterShovel struct {
	shovelMetricsGauge map[string]*prometheus.GaugeVec
}

func newExporterShovel() Exporter {
	return exporterShovel{
		shovelMetricsGauge: shovelGaugeVec,
	}
}

func (e exporterShovel) String() string {
	return "Exporter shovel"
}

func (e exporterShovel) Collect(ch chan<- prometheus.Metric) error {
	shovelData, err := getStatsInfo(config, "shovels", shovelLabelsKeys)

	if err != nil {
		return err
	}

	for key, gauge := range e.shovelMetricsGauge {
		for _, shovel := range shovelData {
			labels := make([]string, len(shovelLabelsKeys))
			val := float64(0)

			for k, v := range shovelLabelsKeys {
				labels[k] = shovel.labels[v]
			}

			if key == "state" && shovel.labels["state"] == "running" {
				val = 1
			} else {
				t, err := time.Parse("2006-01-02 15:04:05", shovel.labels["timestamp"])
				if err != nil {
					return err
				} else {
					val = float64(t.Unix())
				}
			}

			gauge.WithLabelValues(labels...).Set(val)
		}
	}

	for _, gauge := range e.shovelMetricsGauge {
		gauge.Collect(ch)
	}
	return nil
}

func (e exporterShovel) Describe(ch chan<- *prometheus.Desc) {
	for _, shovelMetric := range e.shovelMetricsGauge {
		shovelMetric.Describe(ch)
	}

}
