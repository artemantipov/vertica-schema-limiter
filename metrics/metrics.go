package metrics

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

type Metric struct {
	Value     float64
	CreatedAt time.Time
}

var metricsMap = make(map[string]Metric)

func (metric Metric) IsExpired() (expire bool) {
	expiredAt := metric.CreatedAt.Add(5 * time.Minute)
	expire = expiredAt.Before(time.Now())
	return
}

//Start run metrics http endpoint
func Start() {
	prometheus.MustRegister(newLimiterCollector())
	go checkExpired()
	fmt.Println("Starting metrics server")
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}

//GetMap return metrics map
func GetMap() map[string]Metric {
	return metricsMap
}

func checkExpired() {
	for {
		for metric := range metricsMap {
			if metricsMap[metric].IsExpired() {
				delete(metricsMap, metric)
			}
		}
		time.Sleep(1 * time.Minute)
	}
}
