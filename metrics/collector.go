package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type limiterCollector struct {
	schemaSize *prometheus.Desc
}

func newLimiterCollector() *limiterCollector {
	return &limiterCollector{
		schemaSize: prometheus.NewDesc("vertica_limiter_schema_size",
			"Shows size of checked schemas",
			[]string{"schema"}, nil,
		),
	}
}

func (c *limiterCollector) Describe(ch chan<- *prometheus.Desc) {

	//Update this section with the each metric you create for a given collector
	ch <- c.schemaSize
}

func (c *limiterCollector) Collect(ch chan<- prometheus.Metric) {

	//Implement logic here to determine proper metric value to return to prometheus
	//for each descriptor or call other functions that do so.
	for name, metric := range GetMap() {
		//Write latest value for each metric in the prometheus metric channel.
		//Note that you can pass CounterValue, GaugeValue, or UntypedValue types here.
		ch <- prometheus.MustNewConstMetric(c.schemaSize, prometheus.GaugeValue, metric.Value, name)
	}
}
